package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/williamkoller/cloud-architecture-golang/internal/usr/domain"
	"github.com/williamkoller/cloud-architecture-golang/internal/usr/domain/vo"
)

var (
	ErrAlreadyExists = errors.New("user already exists")
	ErrNotFound      = errors.New("user not found")
)

type UserRepository interface {
	Create(ctx context.Context, u domain.User) error
	GetByEmail(ctx context.Context, email vo.Email) (domain.User, bool, error)
	List(ctx context.Context) ([]domain.User, error)
	Update(ctx context.Context, u domain.User) error
	Delete(ctx context.Context, email vo.Email) error
}


type inMemoryUserRepo struct {
	mu   sync.RWMutex
	data map[string]domain.User
}

func NewInMemoryUserRepository() UserRepository {
	return &inMemoryUserRepo{
		data: make(map[string]domain.User),
	}
}

func (r *inMemoryUserRepo) Create(ctx context.Context, u domain.User) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	key := string(u.Email)

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.data[key]; ok {
		return ErrAlreadyExists
	}
	// cópia defensiva
	r.data[key] = u
	return nil
}

func (r *inMemoryUserRepo) GetByEmail(ctx context.Context, email vo.Email) (domain.User, bool, error) {
	select {
	case <-ctx.Done():
		return domain.User{}, false, ctx.Err()
	default:
	}
	key := string(email)

	r.mu.RLock()
	defer r.mu.RUnlock()

	u, ok := r.data[key]
	return u, ok, nil
}

func (r *inMemoryUserRepo) List(ctx context.Context) ([]domain.User, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]domain.User, 0, len(r.data))
	for _, u := range r.data {
		out = append(out, u)
	}
	return out, nil
}

func (r *inMemoryUserRepo) Update(ctx context.Context, u domain.User) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	key := string(u.Email)

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.data[key]; !ok {
		return ErrNotFound
	}
	r.data[key] = u
	return nil
}

func (r *inMemoryUserRepo) Delete(ctx context.Context, email vo.Email) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	key := string(email)

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.data[key]; !ok {
		return ErrNotFound
	}
	delete(r.data, key)
	return nil
}
