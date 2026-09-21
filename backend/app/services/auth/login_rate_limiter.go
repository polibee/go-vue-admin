package authservices

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/goravel/framework/contracts/database/orm"

	"goravel/app/facades"
)

const (
	loginRateLimitMaxAttempts = 5
	loginRateLimitWindow      = 10 * time.Minute
)

var ErrLoginRateLimitStoreUnavailable = errors.New("login rate limit store unavailable")

type LoginRateLimiter struct{}

func NewLoginRateLimiter() *LoginRateLimiter { return &LoginRateLimiter{} }

func loginRateLimitKey(email, ip string) string {
	identity := strings.ToLower(strings.TrimSpace(email)) + "\x00" + strings.TrimSpace(ip)
	digest := sha256.Sum256([]byte(identity))
	return hex.EncodeToString(digest[:])
}

func (s *LoginRateLimiter) Allow(email, ip string) (bool, error) {
	key := "auth:login:attempts:" + loginRateLimitKey(email, ip)
	cache := facades.Cache()
	cache.Add(key, int64(0), loginRateLimitWindow)
	count, err := cache.Increment(key, 1)
	if err == nil {
		return count <= loginRateLimitMaxAttempts, nil
	}

	return s.allowPostgres(loginRateLimitKey(email, ip), time.Now())
}

func (s *LoginRateLimiter) Reset(email, ip string) {
	keyHash := loginRateLimitKey(email, ip)
	_ = facades.Cache().Forget("auth:login:attempts:" + keyHash)
	_, _ = facades.Orm().Query().Table("auth_login_attempts").Where("key_hash = ?", keyHash).Delete()
}

func (s *LoginRateLimiter) allowPostgres(keyHash string, now time.Time) (bool, error) {
	allowed := false
	err := facades.Orm().Transaction(func(tx orm.Query) error {
		var row struct {
			KeyHash         string    `gorm:"column:key_hash"`
			Attempts        int       `gorm:"column:attempts"`
			WindowStartedAt time.Time `gorm:"column:window_started_at"`
		}
		query := tx.Table("auth_login_attempts").Where("key_hash = ?", keyHash).LockForUpdate()
		if err := query.First(&row); err != nil {
			if createErr := tx.Table("auth_login_attempts").Create(&map[string]any{
				"key_hash": keyHash, "attempts": 1, "window_started_at": now, "updated_at": now,
			}); createErr != nil {
				return createErr
			}
			allowed = true
			return nil
		}

		if now.Sub(row.WindowStartedAt) >= loginRateLimitWindow {
			row.Attempts = 1
			row.WindowStartedAt = now
		} else {
			row.Attempts++
		}
		allowed = row.Attempts <= loginRateLimitMaxAttempts
		_, err := tx.Table("auth_login_attempts").Where("key_hash = ?", keyHash).Update(map[string]any{
			"attempts": row.Attempts, "window_started_at": row.WindowStartedAt, "updated_at": now,
		})
		return err
	})
	if err != nil {
		return false, fmt.Errorf("%w: %v", ErrLoginRateLimitStoreUnavailable, err)
	}
	return allowed, nil
}
