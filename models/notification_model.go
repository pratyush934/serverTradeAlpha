package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/pratyush934/tradealpha/server/database"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type NotificationModel struct {
	Id         string    `gorm:"primaryId; type:varchar(151)" json:"id"`
	UserId     string    `gorm:"userId" json:"userId"`
	Message    string    `gorm:"userId" json:"message"`
	ReadStatus bool      `gorm:"default:false" json:"readStatus"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

func (n *NotificationModel) BeforeCreate(tx *gorm.DB) error {
	n.Id = uuid.New().String()
	n.CreatedAt = time.Now()
	n.UpdatedAt = time.Now()

	return nil
}

func (n *NotificationModel) BeforeUpdate(tx *gorm.DB) error {
	n.UpdatedAt = time.Now()
	return nil
}

func (n *NotificationModel) CreateNotification() (*NotificationModel, error) {
	if err := database.DB.Create(n).Error; err != nil {
		log.Error().Err(err).Msg("issue lies at notification_model/CreateNotification")
		return nil, err
	}
	return n, nil
}

func GetNotificationByUserId(userId string) (*[]NotificationModel, error) {
	var notification []NotificationModel
	if err := database.DB.Where("user_id = ?", userId).Find(&notification).Error; err != nil {
		log.Error().Err(err).Msg("issue lies at notification_model/GetNotificationByUserId")
		return nil, err
	}
	return &notification, nil
}

func UpdateNotification(notification *NotificationModel) error {
	return database.DB.Updates(notification).Error
}

func GetNotificationByNotificationId(id string) (*NotificationModel, error) {
	var notification NotificationModel
	if err := database.DB.Where("notification_id = ?", id).First(&notification).Error; err != nil {
		log.Error().Err(err).Msg("issue lies at notification_model/GetNotificationByNotificationId")
		return nil, err
	}
	return &notification, nil
}

func DeleteNotification(userId string) error {
	return database.DB.Where("user_id = ?", userId).Error
}

func DeleteNotificationByNId(id string) error {
	return database.DB.Where("notification_id = ?", id).Error
}

func MarkAsRead(id string) error {
	ByNotificationId, err := GetNotificationByNotificationId(id)
	if err != nil {
		log.Error().Err(err).Msg("issue lies in the notification_model/MarkAsRead")
		return err
	}
	ByNotificationId.ReadStatus = true
	return nil
}

func GetUnreadNotificationsCount() (int64, error) {
	var count int64

	if err := database.DB.Model(&NotificationModel{}).Where("read_status = ?", false).Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("issue lies in the notification_model/GetUnreadNotification")
		return 0, err
	}
	return count, nil
}
