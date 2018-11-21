package usersvc

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"context"

	"github.com/gorilla/mux"

	httptransport "github.com/go-kit/kit/transport/http"
)

var (
	// ErrBadRouting is returned when an expected path variable is missing.
	// It always indicates programmer error.
	ErrBadRouting = errors.New("inconsistent mapping between route and handler (programmer error)")
)

func MakeHttpHandler(s Service) http.Handler {
	r := mux.NewRouter()
	e := MakeServerEndpoints(s)

	// GET /users/:id     fetchs a user
	// PUT /users/:id     creates a user
	// DELETE /users/:id  deletes a user

	r.Methods("GET").Path("/users/{id}").Handler(httptransport.NewServer(
		e.GetUserEndpoint,
		decodeGetUserRequest,
		encodeResponse,
	))

	r.Methods("PUT").Path("/users/{id}").Handler(httptransport.NewServer(
		e.PutUserEndpoint,
		decodePutUserRequest,
		encodeResponse,
	))

	r.Methods("DELETE").Path("/users/{id}").Handler(httptransport.NewServer(
		e.DeleteUserEndpoint,
		decodeDeleteUserRequest,
		encodeResponse,
	))

	return r

}

func decodeGetUserRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	fmt.Printf("GET\n")

	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}

	r.ParseForm()
	password := r.Form.Get("password")

	return getUserRequest{ID: id, Password: password}, nil
}

func decodePutUserRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	fmt.Printf("PUT\n")

	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}

	var user User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		fmt.Println(r.Body)
		return nil, err
	}
	return putUserRequest{
		ID:   id,
		User: user,
	}, nil
}

func decodeDeleteUserRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	fmt.Printf("DELETE\n")

	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}

	var user User

	if e := json.NewDecoder(r.Body).Decode(&user); e != nil {
		return nil, e
	}
	return deleteUserRequest{
		ID:   id,
		User: user,
	}, nil
}

type errorer interface {
	error() error
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	case ErrNotFound:
		return http.StatusNotFound
	case ErrInconsistentIDs:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
