package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/pratyush934/tradealpha/server/database"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type User struct {
	Id                 string              `gorm:"primaryKey;type:varchar(151)" json:"id"`
	Name               string              `gorm:"not null" json:"name"`
	Email              string              `gorm:"not null;unique" json:"email"`
	PhoneNumber        string              `json:"phoneNumber"`
	ProfileImage       string              `json:"profileImage"`
	AccountBalance     float64             `gorm:"default:0" json:"accountBalance"`
	RoleId             int                 `gorm:"not null;default:1" json:"roleId"`
	Address            []AddressModel      `gorm:"foreignKey:userId" json:"address"`
	PortFolio          []PortFolio         `gorm:"foreignKey:userId" json:"portFolio"`
	Transactions       []TransactionModel  `gorm:"foreignKey:userId" json:"transactions"`
	Notification       []NotificationModel `gorm:"foreignKey:userId" json:"notification"`
	Role               Role                `gorm:"not null;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"role"`
	VerificationStatus bool                `gorm:"default:false" json:"verificationStatus"`
	IsActive           bool                `json:"isActive"`
	Referral           string              `json:"referral"`
	LastLogin          time.Time           `json:"lastLogin"`
	CreatedAt          time.Time           `json:"createdAt"`
	UpdatedAt          time.Time           `json:"updatedAt"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.Id = uuid.New().String()
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()

	return nil
}

func (u *User) CreateUser() (*User, error) {
	if err := database.DB.Create(u).Error; err != nil {
		log.Error().Err(err).Msg("Issue lie in the user_model/CreateUser")
		return nil, err
	}
	return u, nil
}

func GetUserById(id string) (*User, error) {
	var user User
	if err := database.DB.Where("user_id = ?").First(&user).Error; err != nil {
		log.Error().Err(err).Msg("Issue lie in the user_model/GetUserById")
		return nil, err
	}
	return &user, nil
}

func GetUserByEmail(id string) (*User, error) {
	var user User
	if err := database.DB.Where("email = ?").First(&user).Error; err != nil {
		log.Error().Err(err).Msg("Issue lie in the user_model/GetUserByEmail")
		return nil, err
	}
	return &user, nil
}

func UpdateUser(user User) error {
	if err := database.DB.Updates(&user).Error; err != nil {
		log.Error().Err(err).Msg("Issue lie in the user_model/UpdateUser")
		return err
	}
	return nil
}

func GetAllUsers() ([]User, error) {
	var users []User
	if err := database.DB.Find(&users).Error; err != nil {
		log.Error().Err(err).Msg("Issue lie the user_model/GetAllUsers")
		return nil, err
	}
	return users, nil
}

func DeleteUserById(id string) error {
	return database.DB.Where("user_id = ?", id).Delete(&User{}).Error
}

func DeleteUserByEmail(email string) error {
	return database.DB.Where("email = ?", email).Delete(&User{}).Error
}
