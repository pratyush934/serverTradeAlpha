package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/pratyush934/tradealpha/server/database"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

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

func (a *AddressModel) BeforeCreate(tx *gorm.DB) error {
	a.Id = uuid.New().String()
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()

	return nil
}

func (a *AddressModel) BeforeUpdate(tx *gorm.DB) error {
	a.UpdatedAt = time.Now()
	return nil
}

func (a *AddressModel) CreateAddress() (*AddressModel, error) {

	if err := database.DB.Create(a).Error; err != nil {
		log.Error().Err(err).Msg("issue lie in address_model/CreateAddress")
		return nil, err
	}
	return a, nil
}

func GetAddressByUserId(userId string) (*[]AddressModel, error) {
	var address []AddressModel
	if err := database.DB.Where("user_id = ?", userId).Find(&address).Error; err != nil {
		log.Error().Err(err).Msg("issue lie in address_model/AddressByUserId")
		return nil, err
	}
	return &address, nil
}

func GetAddressByAddressId(id string) (*AddressModel, error) {
	var address AddressModel
	if err := database.DB.Where("id = ?", id).Find(&address).Error; err != nil {
		log.Error().Err(err).Msg("issue lie in address_model/GetAddressByAddressId")
		return nil, err
	}
	return &address, nil
}

func UpdateAddress(address *AddressModel) (*AddressModel, error) {
	if err := database.DB.Updates(address).Error; err != nil {
		log.Error().Err(err).Msg("issue lie in address_model/UpdateAddress")
		return nil, err
	}
	return address, nil
}

func DeleteAddress(userId string) error {
	return database.DB.Where("user_id = ?", userId).Delete(&AddressModel{}).Error
}

func DeleteAddressByAddressId(id string) error {
	return database.DB.Where("id = ?", id).Delete(&AddressModel{}).Error
}
