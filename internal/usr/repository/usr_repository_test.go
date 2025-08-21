package repository

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/williamkoller/cloud-architecture-golang/internal/usr/domain"
	"github.com/williamkoller/cloud-architecture-golang/internal/usr/domain/vo"
)


func mustEmail(t *testing.T, s string) vo.Email {
	t.Helper()
	e, err := vo.NewEmail(s)
	if err != nil {
		t.Fatalf("invalid email %q: %v", s, err)
	}
	return e
}

func mustUser(t *testing.T, name, email string, active bool, ut domain.UserType) domain.User {
	t.Helper()
	u, err := domain.NewUser(name, email, "secret123", active, ut)
	if err != nil {
		t.Fatalf("NewUser: %v", err)
	}
	return u
}


func TestCreateAndGetByEmail(t *testing.T) {
	repo := NewInMemoryUserRepository()

	u := mustUser(t, "Ana", "ana@example.com", true, domain.UserTypeUser)

	if err := repo.Create(context.Background(), u); err != nil {
		t.Fatalf("Create: %v", err)
	}

	got, ok, err := repo.GetByEmail(context.Background(), mustEmail(t, "ana@example.com"))
	if err != nil {
		t.Fatalf("GetByEmail: %v", err)
	}
	if !ok {
		t.Fatalf("GetByEmail: expected ok=true")
	}

	if got.Name != u.Name || got.Active != u.Active || got.UserType != u.UserType || string(got.Email) != string(u.Email) {
		t.Fatalf("GetByEmail: got %+v, want %+v", got, u)
	}
}

func TestCreate_DuplicateConcurrent(t *testing.T) {
	repo := NewInMemoryUserRepository()
	u := mustUser(t, "Ana", "ana@example.com", true, domain.UserTypeUser)

	const N = 16
	var wg sync.WaitGroup
	var successes int64
	var conflicts int64

	wg.Add(N)
	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			err := repo.Create(context.Background(), u)
			switch err {
			case nil:
				atomic.AddInt64(&successes, 1)
			case ErrAlreadyExists:
				atomic.AddInt64(&conflicts, 1)
			default:
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		}()
	}
	wg.Wait()

	if successes != 1 {
		t.Fatalf("expected exactly 1 success, got %d", successes)
	}
	if conflicts != N-1 {
		t.Fatalf("expected %d conflicts, got %d", N-1, conflicts)
	}

	_, ok, err := repo.GetByEmail(context.Background(), mustEmail(t, "ana@example.com"))
	if err != nil || !ok {
		t.Fatalf("GetByEmail after concurrent create: ok=%v err=%v", ok, err)
	}
}

func TestGetByEmail_NotFound(t *testing.T) {
	repo := NewInMemoryUserRepository()

	_, ok, err := repo.GetByEmail(context.Background(), mustEmail(t, "missing@example.com"))
	if err != nil {
		t.Fatalf("GetByEmail: unexpected err: %v", err)
	}
	if ok {
		t.Fatalf("GetByEmail: expected ok=false for missing user")
	}
}

func TestList_ReturnsCopies(t *testing.T) {
	repo := NewInMemoryUserRepository()

	u1 := mustUser(t, "Ana", "ana@example.com", true, domain.UserTypeUser)
	u2 := mustUser(t, "Bob", "bob@example.com", false, domain.UserTypeAdmin)

	if err := repo.Create(context.Background(), u1); err != nil {
		t.Fatalf("Create u1: %v", err)
	}
	if err := repo.Create(context.Background(), u2); err != nil {
		t.Fatalf("Create u2: %v", err)
	}

	list, err := repo.List(context.Background())
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("List length: got %d, want 2", len(list))
	}

	list[0].Name = "Hacked"
	got, ok, err := repo.GetByEmail(context.Background(), mustEmail(t, "ana@example.com"))
	if err != nil || !ok {
		t.Fatalf("GetByEmail ana: ok=%v err=%v", ok, err)
	}
	if got.Name != "Ana" {
		t.Fatalf("repository was mutated via List result; got.Name=%q", got.Name)
	}
}

func TestUpdate_Success_And_NotFound(t *testing.T) {
	repo := NewInMemoryUserRepository()

	u := mustUser(t, "Ana", "ana@example.com", true, domain.UserTypeUser)
	if err := repo.Create(context.Background(), u); err != nil {
		t.Fatalf("Create: %v", err)
	}

	updated := mustUser(t, "Ana Paula", "ana@example.com", false, domain.UserTypeAdmin)
	if err := repo.Update(context.Background(), updated); err != nil {
		t.Fatalf("Update existing: %v", err)
	}

	got, ok, err := repo.GetByEmail(context.Background(), mustEmail(t, "ana@example.com"))
	if err != nil || !ok {
		t.Fatalf("GetByEmail after update: ok=%v err=%v", ok, err)
	}
	if got.Name != "Ana Paula" || got.Active != false || got.UserType != domain.UserTypeAdmin {
		t.Fatalf("updated user mismatch: got=%+v", got)
	}

	missing := mustUser(t, "Carlos", "carlos@example.com", true, domain.UserTypeUser)
	if err := repo.Update(context.Background(), missing); err != ErrNotFound {
		t.Fatalf("Update missing: got %v, want %v", err, ErrNotFound)
	}
}

func TestDelete_Success_And_NotFound(t *testing.T) {
	repo := NewInMemoryUserRepository()

	u := mustUser(t, "Ana", "ana@example.com", true, domain.UserTypeUser)
	if err := repo.Create(context.Background(), u); err != nil {
		t.Fatalf("Create: %v", err)
	}

	if err := repo.Delete(context.Background(), mustEmail(t, "ana@example.com")); err != nil {
		t.Fatalf("Delete existing: %v", err)
	}

	_, ok, err := repo.GetByEmail(context.Background(), mustEmail(t, "ana@example.com"))
	if err != nil {
		t.Fatalf("GetByEmail after delete: %v", err)
	}
	if ok {
		t.Fatalf("GetByEmail after delete: expected ok=false")
	}

	if err := repo.Delete(context.Background(), mustEmail(t, "ana@example.com")); err != ErrNotFound {
		t.Fatalf("Delete missing: got %v, want %v", err, ErrNotFound)
	}
}

func TestContextCanceled_OnOperations(t *testing.T) {
	repo := NewInMemoryUserRepository()
	u := mustUser(t, "Ana", "ana@example.com", true, domain.UserTypeUser)

	ctxC, cancel := context.WithCancel(context.Background())
	cancel()
	if err := repo.Create(ctxC, u); err == nil {
		t.Fatalf("Create with canceled context: expected error")
	}

	ctxG, cancelG := context.WithCancel(context.Background())
	cancelG()
	if _, _, err := repo.GetByEmail(ctxG, mustEmail(t, "ana@example.com")); err == nil {
		t.Fatalf("GetByEmail with canceled context: expected error")
	}

	ctxL, cancelL := context.WithCancel(context.Background())
	cancelL()
	if _, err := repo.List(ctxL); err == nil {
		t.Fatalf("List with canceled context: expected error")
	}

	ctxU, cancelU := context.WithCancel(context.Background())
	cancelU()
	if err := repo.Update(ctxU, u); err == nil {
		t.Fatalf("Update with canceled context: expected error")
	}

	ctxD, cancelD := context.WithCancel(context.Background())
	cancelD()
	if err := repo.Delete(ctxD, mustEmail(t, "ana@example.com")); err == nil {
		t.Fatalf("Delete with canceled context: expected error")
	}
}
