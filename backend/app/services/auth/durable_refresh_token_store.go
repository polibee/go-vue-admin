package authservices

import (
	"strconv"
	"time"

	"goravel/app/facades"
)

type cacheRefreshTokenStore struct {
	cache interface {
		Put(key string, value any, ttl time.Duration) error
		GetString(key string, def ...string) string
		Forget(key string) bool
	}
}

func (s cacheRefreshTokenStore) Put(key string, value any, ttl time.Duration) error {
	return s.cache.Put(key, value, ttl)
}

func (s cacheRefreshTokenStore) GetString(key string, def ...string) string {
	return s.cache.GetString(key, def...)
}

func (s cacheRefreshTokenStore) Forget(key string) bool {
	return s.cache.Forget(key)
}

func (s cacheRefreshTokenStore) Probe() error {
	const key = "auth:refresh:availability-probe"
	if err := s.cache.Put(key, "ok", time.Second); err != nil {
		return err
	}
	s.cache.Forget(key)
	return nil
}

func (cacheRefreshTokenStore) RevokeAll(_ uint) error { return nil }

type postgresRefreshTokenStore struct{}

func (postgresRefreshTokenStore) Put(key string, value any, ttl time.Duration) error {
	userID, err := strconv.ParseUint(value.(string), 10, 64)
	if err != nil {
		return err
	}
	return facades.Orm().Query().Table("auth_refresh_tokens").Create(&map[string]any{
		"token_hash": key,
		"user_id":    userID,
		"expires_at": time.Now().Add(ttl),
		"created_at": time.Now(),
		"updated_at": time.Now(),
	})
}

func (postgresRefreshTokenStore) GetString(key string, _ ...string) string {
	var row struct {
		UserID    uint      `gorm:"column:user_id"`
		ExpiresAt time.Time `gorm:"column:expires_at"`
	}
	if err := facades.Orm().Query().Table("auth_refresh_tokens").Where("token_hash = ? AND expires_at > ?", key, time.Now()).First(&row); err != nil {
		return ""
	}
	return strconv.FormatUint(uint64(row.UserID), 10)
}

func (postgresRefreshTokenStore) Forget(key string) bool {
	_, err := facades.Orm().Query().Table("auth_refresh_tokens").Where("token_hash = ?", key).Delete()
	return err == nil
}

func (postgresRefreshTokenStore) Probe() error {
	_, err := facades.Orm().Query().Table("auth_refresh_tokens").Count()
	return err
}

func (postgresRefreshTokenStore) RevokeAll(userID uint) error {
	_, err := facades.Orm().Query().Table("auth_refresh_tokens").Where("user_id = ?", userID).Delete()
	return err
}

type durableRefreshTokenStore struct {
	database refreshTokenStore
	cache    refreshTokenStore
}

func (s durableRefreshTokenStore) Put(key string, value any, ttl time.Duration) error {
	if err := s.database.Put(key, value, ttl); err != nil {
		return err
	}
	_ = s.cache.Put(key, value, ttl)
	return nil
}

func (s durableRefreshTokenStore) GetString(key string, def ...string) string {
	// PostgreSQL is authoritative so logout-all cannot be bypassed by a stale Redis mirror.
	return s.database.GetString(key, def...)
}

func (s durableRefreshTokenStore) Forget(key string) bool {
	databaseOK := s.database.Forget(key)
	_ = s.cache.Forget(key)
	return databaseOK
}

func (s durableRefreshTokenStore) Probe() error {
	return s.database.Probe()
}

func (s durableRefreshTokenStore) RevokeAll(userID uint) error {
	return s.database.RevokeAll(userID)
}

func NewDurableRefreshTokenService() *RefreshTokenService {
	return NewRefreshTokenService(durableRefreshTokenStore{
		database: postgresRefreshTokenStore{},
		cache:    cacheRefreshTokenStore{cache: facades.Cache()},
	})
}
