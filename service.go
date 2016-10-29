package usersvc

import (
	"errors"
	"fmt"
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
	mtx sync.RWMutex
	m   map[string]User
}

func NewUserService() Service {
	return &userService{
		m: map[string]User{},
	}
}

func (s *userService) GetUser(ctx context.Context, id string, password string) (User, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	u, ok := s.m[id]

	if !ok {
		return User{}, ErrNotFound
	}

	return u, nil

}

func (s *userService) PutUser(ctx context.Context, id string, u User) (User, error) {
	if id != u.ID {
		return User{}, ErrInconsistentIDs
	}

	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.m[id] = u
	return u, nil
}

func (s *userService) DeleteUser(ctx context.Context, id string, u User) (User, error) {
	fmt.Printf(id)
	fmt.Printf("\n")
	fmt.Printf(u.ID)
	fmt.Printf("\n")
	if id != u.ID {
		return User{}, ErrInconsistentIDs
	}

	s.mtx.Lock()
	defer s.mtx.Unlock()
	if _, ok := s.m[id]; !ok {
		return User{}, ErrNotFound
	}
	delete(s.m, id)
	return u, nil
}
