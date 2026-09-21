package services

import (
	"errors"
	"testing"
	"time"
)

type fakeRefreshTokenStore struct {
	values map[string]string
	fail   bool
}

func (f *fakeRefreshTokenStore) Put(key string, value any, _ time.Duration) error {
	if f.fail {
		return errors.New("cache offline")
	}
	f.values[key] = value.(string)
	return nil
}

func (f *fakeRefreshTokenStore) GetString(key string, _ ...string) string {
	return f.values[key]
}

func (f *fakeRefreshTokenStore) Forget(key string) bool {
	delete(f.values, key)
	return true
}

func (f *fakeRefreshTokenStore) Probe() error {
	if f.fail {
		return errors.New("cache offline")
	}
	return nil
}

func (f *fakeRefreshTokenStore) RevokeAll(_ uint) error { return nil }

func TestRefreshTokenIssueAndConsumeRotatesOpaqueToken(t *testing.T) {
	store := &fakeRefreshTokenStore{values: map[string]string{}}
	service := NewRefreshTokenService(store)

	first, err := service.Issue(42)
	if err != nil {
		t.Fatalf("issue refresh token: %v", err)
	}
	if first == "" {
		t.Fatal("expected an opaque refresh token")
	}

	userID, err := service.Consume(first)
	if err != nil {
		t.Fatalf("consume refresh token: %v", err)
	}
	if userID != 42 {
		t.Fatalf("expected user 42, got %d", userID)
	}
	if _, err := service.Consume(first); err == nil {
		t.Fatal("expected a consumed token to be rejected on replay")
	}
}

func TestRefreshTokenStoreFailureIsReportedSeparately(t *testing.T) {
	service := NewRefreshTokenService(&fakeRefreshTokenStore{values: map[string]string{}, fail: true})

	if _, err := service.Issue(42); !errors.Is(err, ErrRefreshTokenStoreUnavailable) {
		t.Fatalf("expected store unavailable on issue, got %v", err)
	}
	if _, err := service.Consume("opaque-token"); !errors.Is(err, ErrRefreshTokenStoreUnavailable) {
		t.Fatalf("expected store unavailable on consume, got %v", err)
	}
}

func TestDurableStoreFallsBackToPostgresWhenRedisProbeFails(t *testing.T) {
	database := &fakeRefreshTokenStore{values: map[string]string{"auth:refresh:hash": "42"}}
	redis := &fakeRefreshTokenStore{values: map[string]string{}, fail: true}
	store := durableRefreshTokenStore{database: database, cache: redis}

	if got := store.GetString("auth:refresh:hash"); got != "42" {
		t.Fatalf("expected PostgreSQL fallback user 42, got %q", got)
	}
}
