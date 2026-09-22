package notifications

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"goravel/app/facades"
	"goravel/app/models"
)

var (
	ErrInvalidNotification  = errors.New("invalid notification")
	ErrNotificationNotFound = errors.New("notification not found")
)

type NotificationInput struct {
	Type  string
	Title string
	Body  string
	URL   string
}

type NotificationPage struct {
	Data     []models.Notification `json:"data"`
	Page     int                   `json:"page"`
	PerPage  int                   `json:"per_page"`
	Total    int64                 `json:"total"`
	LastPage int                   `json:"last_page"`
}

func NewNotificationService() *NotificationService { return &NotificationService{} }

type NotificationService struct{}

func validateNotificationURL(value string) error {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	if !strings.HasPrefix(value, "/") || strings.HasPrefix(value, "//") || strings.ContainsAny(value, "\r\n") {
		return ErrInvalidNotification
	}
	return nil
}

func validateNotificationInput(userID uint, input NotificationInput) error {
	if userID == 0 || strings.TrimSpace(input.Type) == "" || strings.TrimSpace(input.Title) == "" || strings.TrimSpace(input.Body) == "" {
		return ErrInvalidNotification
	}
	return validateNotificationURL(input.URL)
}

func (s *NotificationService) Create(userID uint, input NotificationInput) (*models.Notification, error) {
	if err := validateNotificationInput(userID, input); err != nil {
		return nil, err
	}
	notification := &models.Notification{
		UserID: userID,
		Type:   strings.TrimSpace(input.Type),
		Title:  strings.TrimSpace(input.Title),
		Body:   strings.TrimSpace(input.Body),
		URL:    strings.TrimSpace(input.URL),
	}
	if err := facades.Orm().Query().Create(notification); err != nil {
		return nil, err
	}
	return notification, nil
}

func normalizePage(page, perPage int) (int, int) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	if perPage > 100 {
		perPage = 100
	}
	return page, perPage
}

func (s *NotificationService) List(userID uint, page, perPage int, unreadOnly bool) (NotificationPage, error) {
	if userID == 0 {
		return NotificationPage{}, ErrInvalidNotification
	}
	page, perPage = normalizePage(page, perPage)
	query := facades.Orm().Query().Where("user_id = ?", userID)
	if unreadOnly {
		query = query.WhereNull("read_at")
	}
	var data []models.Notification
	var total int64
	if err := query.OrderByDesc("id").Paginate(page, perPage, &data, &total); err != nil {
		return NotificationPage{}, err
	}
	lastPage := 1
	if total > 0 {
		lastPage = int((total + int64(perPage) - 1) / int64(perPage))
	}
	return NotificationPage{Data: data, Page: page, PerPage: perPage, Total: total, LastPage: lastPage}, nil
}

func (s *NotificationService) UnreadCount(userID uint) (int64, error) {
	if userID == 0 {
		return 0, ErrInvalidNotification
	}
	return facades.Orm().Query().Where("user_id = ?", userID).WhereNull("read_at").Count()
}

func (s *NotificationService) MarkRead(userID, id uint) error {
	if userID == 0 || id == 0 {
		return ErrInvalidNotification
	}
	var notification models.Notification
	if err := facades.Orm().Query().Where("user_id = ? AND id = ?", userID, id).First(&notification); err != nil {
		return ErrNotificationNotFound
	}
	if notification.ReadAt != nil {
		return nil
	}
	now := time.Now()
	result, err := facades.Orm().Query().Model(&models.Notification{}).Where("user_id = ? AND id = ? AND read_at IS NULL", userID, id).Update(map[string]any{"read_at": now, "updated_at": now})
	if err != nil {
		return err
	}
	if result.RowsAffected == 0 {
		return nil
	}
	return nil
}

func (s *NotificationService) MarkAllRead(userID uint) (int64, error) {
	if userID == 0 {
		return 0, ErrInvalidNotification
	}
	now := time.Now()
	result, err := facades.Orm().Query().Model(&models.Notification{}).Where("user_id = ? AND read_at IS NULL", userID).Update(map[string]any{"read_at": now, "updated_at": now})
	if err != nil {
		return 0, err
	}
	return result.RowsAffected, nil
}

func notificationError(code string, err error) error {
	return fmt.Errorf("%s: %w", code, err)
}
