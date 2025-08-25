package models

import "time"

type AddressModel struct {
	Id        string    `gorm:"primaryKey" json:"id"`
	UserId    string    `gorm:"not null" json:"userId"`
	Street    string    `gorm:"not null" json:"street"`
	ZipCode   string    `gorm:"not null" json:"zipCode"`
	City      string    `gorm:"not null" json:"city"`
	State     string    `gorm:"not null" json:"state"`
	Country   string    `gorm:"not null" json:"country"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
