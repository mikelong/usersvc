package usersvc

import (
	"errors"
	"sync"

	"golang.org/x/net/context"
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
)

type userService struct {
	r   Repository
	mtx sync.RWMutex
	m   map[string]User
}

func NewUserService() Service {
	return &userService{
		r: NewRepository(),
		m: map[string]User{},
	}
}

func (s *userService) GetUser(ctx context.Context, id string, password string) (User, error) {
	u, err := s.r.GetUser(id, password)
	return u, err
}

func (s *userService) PutUser(ctx context.Context, id string, u User) (User, error) {
	if id != u.ID {
		return User{}, ErrInconsistentIDs
	}

	err := s.r.PutUser(u)
	return u, err
}

func (s *userService) DeleteUser(ctx context.Context, id string, u User) (User, error) {
	if id != u.ID {
		return User{}, ErrInconsistentIDs
	}

	err := s.r.DeleteUser(u)
	return u, err
}
