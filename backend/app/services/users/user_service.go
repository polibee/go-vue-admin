package userservices

import (
	"errors"
	"fmt"
	"strings"

	"github.com/goravel/framework/contracts/database/orm"

	"goravel/app/facades"
	"goravel/app/models"
	notificationservices "goravel/app/services/notifications"
	rbacservices "goravel/app/services/rbac"
)

var (
	ErrUserExists  = errors.New("user already exists")
	ErrInvalidUser = errors.New("invalid user")
)

const (
	userStatusActive   = "active"
	userStatusDisabled = "disabled"
	userStatusLocked   = "locked"
)

type UserService struct{}

func NewUserService() *UserService { return &UserService{} }

func validateUserInput(name, email, password string, requirePassword bool) error {
	if strings.TrimSpace(name) == "" || strings.TrimSpace(email) == "" {
		return ErrInvalidUser
	}
	if requirePassword && strings.TrimSpace(password) == "" {
		return ErrInvalidUser
	}
	if strings.TrimSpace(password) != "" && len(strings.TrimSpace(password)) < 8 {
		return ErrInvalidUser
	}
	return nil
}

func validateUserStatus(status string) error {
	switch strings.TrimSpace(status) {
	case userStatusActive, userStatusDisabled, userStatusLocked:
		return nil
	default:
		return ErrInvalidUser
	}
}

func normalizeUserStatus(status string) string {
	if strings.TrimSpace(status) == "" {
		return userStatusActive
	}
	return strings.TrimSpace(status)
}

func (s *UserService) Create(name, email, password, locale, status string) (*models.User, error) {
	name, email, password = strings.TrimSpace(name), strings.TrimSpace(email), strings.TrimSpace(password)
	if err := validateUserInput(name, email, password, true); err != nil {
		return nil, err
	}
	status = normalizeUserStatus(status)
	if err := validateUserStatus(status); err != nil {
		return nil, err
	}
	exists, err := facades.Orm().Query().Where("email = ?", email).Exists()
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUserExists
	}
	hash, err := facades.Hash().Make(password)
	if err != nil {
		return nil, err
	}
	if locale == "" {
		locale = "zh-CN"
	}
	user := &models.User{Name: name, Email: email, Password: hash, Locale: locale, Status: status}
	if err := facades.Orm().Query().Create(user); err != nil {
		return nil, err
	}
	notificationservices.NewNotificationService().PublishBestEffort(user.ID, notificationservices.NotificationInput{
		Type:  notificationservices.TypeUserCreated,
		Title: "账号已创建",
		Body:  fmt.Sprintf("管理员已为你创建账号：%s。", user.Name),
		URL:   fmt.Sprintf("/admin/users/%d/edit", user.ID),
	})
	return user, nil
}

func (s *UserService) Update(id int64, name, email, password, locale, status string) (*models.User, error) {
	var user models.User
	if err := facades.Orm().Query().Where("id", id).First(&user); err != nil {
		return nil, ErrUserNotFound
	}
	previousStatus := user.Status
	name, email, password = strings.TrimSpace(name), strings.TrimSpace(email), strings.TrimSpace(password)
	if err := validateUserInput(name, email, password, false); err != nil {
		return nil, err
	}
	if err := validateUserStatus(status); err != nil {
		return nil, err
	}
	exists, err := facades.Orm().Query().Where("email = ? AND id <> ?", email, id).Exists()
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUserExists
	}
	if locale == "" {
		locale = "zh-CN"
	}
	user.Name, user.Email, user.Locale, user.Status = name, email, locale, status
	if password != "" {
		hash, err := facades.Hash().Make(password)
		if err != nil {
			return nil, err
		}
		user.Password = hash
	}
	if err := facades.Orm().Query().Save(&user); err != nil {
		return nil, err
	}
	notifications := notificationservices.NewNotificationService()
	if previousStatus != user.Status {
		notifications.PublishBestEffort(user.ID, notificationservices.NotificationInput{
			Type:  notificationservices.TypeUserStatusChanged,
			Title: "账号状态已变更",
			Body:  fmt.Sprintf("你的账号状态已变更为：%s。", user.Status),
			URL:   fmt.Sprintf("/admin/users/%d/edit", user.ID),
		})
	}
	if password != "" {
		notifications.PublishBestEffort(user.ID, notificationservices.NotificationInput{
			Type:  notificationservices.TypeUserPasswordReset,
			Title: "密码已重置",
			Body:  "管理员已重置你的登录密码，请使用新密码登录。",
			URL:   "/admin/users",
		})
	}
	return &user, nil
}

func (s *UserService) Delete(id int64) error {
	var user models.User
	if err := facades.Orm().Query().Where("id", id).First(&user); err != nil {
		return ErrUserNotFound
	}
	last, err := rbacservices.NewRBACService().IsLastActiveAdmin(id)
	if err != nil {
		return err
	}
	if last {
		return ErrLastAdmin
	}
	return facades.Orm().Transaction(func(tx orm.Query) error {
		if _, err := tx.Table("role_user").Where("user_id = ?", id).Delete(); err != nil {
			return err
		}
		_, err := tx.Table("users").Where("id = ?", id).Delete()
		return err
	})
}

func (s *UserService) BulkSetStatus(ids []int64, status string) error {
	if len(ids) == 0 || validateUserStatus(status) != nil {
		return ErrInvalidUser
	}
	seen := make(map[int64]struct{}, len(ids))
	changed := make(map[uint]struct{}, len(ids))
	for _, id := range ids {
		if id < 1 {
			return ErrInvalidUser
		}
		if _, exists := seen[id]; exists {
			continue
		}
		seen[id] = struct{}{}
		var user models.User
		if err := facades.Orm().Query().Where("id", id).First(&user); err != nil {
			return ErrUserNotFound
		}
		if status != userStatusActive && user.Status == userStatusActive {
			last, err := rbacservices.NewRBACService().IsLastActiveAdmin(id)
			if err != nil {
				return err
			}
			if last {
				return ErrLastAdmin
			}
		}
		if user.Status != status {
			changed[user.ID] = struct{}{}
		}
	}
	if err := facades.Orm().Transaction(func(tx orm.Query) error {
		for id := range seen {
			if _, err := tx.Table("users").Where("id", id).Update("status", status); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}
	for id := range changed {
		notificationservices.NewNotificationService().PublishBestEffort(id, notificationservices.NotificationInput{
			Type:  notificationservices.TypeUserStatusChanged,
			Title: "账号状态已批量变更",
			Body:  fmt.Sprintf("你的账号状态已变更为：%s。", status),
			URL:   fmt.Sprintf("/admin/users/%d/edit", id),
		})
	}
	return nil
}
