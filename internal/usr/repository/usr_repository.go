package repository

import (
	"context"
	"errors"
	"hash/fnv"
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

// shard representa um fragmento do repositório com seu próprio lock
type shard struct {
	mu   sync.RWMutex
	data map[string]domain.User
}

// inMemoryUserRepo implementa repositório em memória com sharding para melhor concorrência
type inMemoryUserRepo struct {
	shards    []*shard
	shardMask uint32
}

// NewInMemoryUserRepository cria um repositório em memória otimizado com sharding
func NewInMemoryUserRepository() UserRepository {
	// Usar menos shards para debugging
	numShards := 4
	shards := make([]*shard, numShards)

	for i := 0; i < numShards; i++ {
		shards[i] = &shard{
			data: make(map[string]domain.User),
		}
	}

	return &inMemoryUserRepo{
		shards:    shards,
		shardMask: uint32(numShards - 1), // Para numShards = 4, mask = 3
	}
}

// getShard retorna o shard apropriado para uma chave
func (r *inMemoryUserRepo) getShard(key string) *shard {
	h := fnv.New32a()
	h.Write([]byte(key))
	return r.shards[h.Sum32()&r.shardMask]
}

func (r *inMemoryUserRepo) Create(ctx context.Context, u domain.User) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	key := string(u.Email)
	shard := r.getShard(key)

	shard.mu.Lock()
	defer shard.mu.Unlock()

	if _, ok := shard.data[key]; ok {
		return ErrAlreadyExists
	}

	// Cópia defensiva
	shard.data[key] = u
	return nil
}

func (r *inMemoryUserRepo) GetByEmail(ctx context.Context, email vo.Email) (domain.User, bool, error) {
	select {
	case <-ctx.Done():
		return domain.User{}, false, ctx.Err()
	default:
	}

	key := string(email)
	shard := r.getShard(key)

	shard.mu.RLock()
	defer shard.mu.RUnlock()

	u, ok := shard.data[key]
	return u, ok, nil
}

func (r *inMemoryUserRepo) List(ctx context.Context) ([]domain.User, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Calcular tamanho total primeiro
	totalSize := 0
	for _, shard := range r.shards {
		shard.mu.RLock()
		totalSize += len(shard.data)
		shard.mu.RUnlock()
	}

	// Pré-alocar slice com tamanho conhecido
	result := make([]domain.User, 0, totalSize)

	// Coletar dados de todos os shards
	for _, shard := range r.shards {
		shard.mu.RLock()
		for _, u := range shard.data {
			result = append(result, u)
		}
		shard.mu.RUnlock()
	}

	return result, nil
}

func (r *inMemoryUserRepo) Update(ctx context.Context, u domain.User) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	key := string(u.Email)
	shard := r.getShard(key)

	shard.mu.Lock()
	defer shard.mu.Unlock()

	if _, ok := shard.data[key]; !ok {
		return ErrNotFound
	}

	shard.data[key] = u
	return nil
}

func (r *inMemoryUserRepo) Delete(ctx context.Context, email vo.Email) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	key := string(email)
	shard := r.getShard(key)

	shard.mu.Lock()
	defer shard.mu.Unlock()

	if _, ok := shard.data[key]; !ok {
		return ErrNotFound
	}

	delete(shard.data, key)
	return nil
}
