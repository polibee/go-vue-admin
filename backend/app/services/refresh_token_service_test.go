package services

import (
	"testing"
	"time"
)

type fakeRefreshTokenStore struct {
	values map[string]string
}

func (f *fakeRefreshTokenStore) Put(key string, value any, _ time.Duration) error {
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
