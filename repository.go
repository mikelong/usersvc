package usersvc

import "sync"

type Repository interface {
	GetUser(id string) (User, error)
	PutUser(u User) error
	DeleteUser(u User) error
}

type inmemRepository struct {
	mtx sync.RWMutex
	m   map[string]User
}

func NewRepository() Repository {
	return &inmemRepository{
		m: map[string]User{},
	}
}

func (r *inmemRepository) GetUser(id string) (User, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	u, ok := r.m[id]

	if !ok {
		return User{}, ErrNotFound
	}

	return u, nil
}

func (r *inmemRepository) PutUser(u User) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.m[u.ID] = u
	return nil
}

func (r *inmemRepository) DeleteUser(u User) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	if _, ok := r.m[u.ID]; !ok {
		return ErrNotFound
	}
	delete(r.m, u.ID)
	return nil
}
