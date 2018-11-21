package usersvc

import (
	"context"
	"errors"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	GetUser(ctx context.Context, id string, password string) (User, error)
	PutUser(ctx context.Context, id string, u User) (User, error)
	DeleteUser(ctx context.Context, id string, u User) (User, error)
}

type User struct {
	ID             string `json:"id"`
	Password       string `json:"password,omitempty"`
	HashedPassword string `json:"hashedPassword,omitempty"`
}

var (
	ErrInconsistentIDs = errors.New("inconsistent IDs")
	ErrNotFound        = errors.New("not found")
	ErrInvalidPassword = errors.New("invalid password")
)

type userService struct {
	r   Repository
	mtx sync.RWMutex
	m   map[string]User
}

func NewUserService() Service {
	return &userService{
		r: NewDynamoRepository(),
		m: map[string]User{},
	}
}

func (s *userService) GetUser(ctx context.Context, id string, password string) (User, error) {
	u, err := s.r.GetUser(id)

	if err != nil {
		return User{}, err
	}

	h := []byte(u.HashedPassword)
	p := []byte(password)

	err = bcrypt.CompareHashAndPassword(h, p)

	if err != nil {
		return User{}, ErrNotFound
	}

	return u, nil
}

func (s *userService) PutUser(ctx context.Context, id string, u User) (User, error) {
	if id != u.ID {
		return User{}, ErrInconsistentIDs
	}

	if validPassword(u.Password) {
		h, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)

		if err != nil {
			return User{}, err
		}

		u.HashedPassword = string(h)
		u.Password = ""
	} else {
		return User{}, ErrInvalidPassword
	}

	err := s.r.PutUser(u)
	return u, err
}

func (s *userService) DeleteUser(ctx context.Context, id string, u User) (User, error) {
	if id != u.ID {
		return User{}, ErrInconsistentIDs
	}

	user, err := s.r.GetUser(id)

	h := []byte(user.HashedPassword)
	p := []byte(u.Password)

	err = bcrypt.CompareHashAndPassword(h, p)

	if err != nil {
		return User{}, ErrNotFound
	}

	err = s.r.DeleteUser(u)

	return user, err
}

func validPassword(password string) bool {
	return password != ""
}
