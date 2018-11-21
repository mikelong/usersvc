package usersvc

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	GetUserEndpoint    endpoint.Endpoint
	PutUserEndpoint    endpoint.Endpoint
	DeleteUserEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		GetUserEndpoint:    MakeGetUserEndpoint(s),
		PutUserEndpoint:    MakePutUserEndpoint(s),
		DeleteUserEndpoint: MakeDeleteUserEndpoint(s),
	}
}

func MakeGetUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getUserRequest)
		u, e := s.GetUser(ctx, req.ID, req.Password)
		return getUserResponse{User: u, Err: e}, nil
	}
}

func MakePutUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(putUserRequest)
		u, e := s.PutUser(ctx, req.ID, req.User)
		return putUserResponse{User: u, Err: e}, nil
	}
}

func MakeDeleteUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(deleteUserRequest)
		u, e := s.DeleteUser(ctx, req.ID, req.User)
		return deleteUserResponse{User: u, Err: e}, nil
	}
}

type getUserRequest struct {
	ID       string
	Password string
}

type putUserRequest struct {
	ID   string
	User User
}

type deleteUserRequest struct {
	ID   string
	User User
}

type getUserResponse struct {
	User User  `json:"user,omitempty"`
	Err  error `json:"err,omitempty"`
}

func (r getUserResponse) error() error { return r.Err }

type putUserResponse struct {
	User User  `json:"user,omitempty"`
	Err  error `json:"err,omitempty"`
}

func (r putUserResponse) error() error { return r.Err }

type deleteUserResponse struct {
	User User  `json:"user,omitempty"`
	Err  error `json:"err,omitempty"`
}

func (r deleteUserResponse) error() error { return r.Err }
