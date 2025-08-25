package models

import "time"

type NotificationModel struct {
	Id         string    `gorm:"primaryId; type:varchar(151)" json:"id"`
	UserId     string    `gorm:"userId" json:"userId"`
	Message    string    `gorm:"userId" json:"message"`
	ReadStatus bool      `gorm:"default:false" json:"readStatus"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}
