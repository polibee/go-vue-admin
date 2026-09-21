package services

import (
	"errors"
	"strings"

	"github.com/goravel/framework/contracts/database/orm"

	"goravel/app/facades"
	"goravel/app/models"
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
	return user, nil
}

func (s *UserService) Update(id int64, name, email, password, locale, status string) (*models.User, error) {
	var user models.User
	if err := facades.Orm().Query().Where("id", id).First(&user); err != nil {
		return nil, ErrUserNotFound
	}
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
	return &user, nil
}

func (s *UserService) Delete(id int64) error {
	var user models.User
	if err := facades.Orm().Query().Where("id", id).First(&user); err != nil {
		return ErrUserNotFound
	}
	last, err := NewRBACService().IsLastActiveAdmin(id)
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
