package models

import "time"

type ProtFolio struct {
	Id        string    `gorm:"primaryKey;type:varchar(151)" json:"id"`
	UserId    string    `gorm:"not null" json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
