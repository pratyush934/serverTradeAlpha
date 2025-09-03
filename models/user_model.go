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
	OAuthId            string              `json:"oauth_id"`
	Provider           string              `json:"provider"`
	Name               string              `gorm:"not null" json:"name"`
	Email              string              `gorm:"not null;unique" json:"email"`
	PhoneNumber        string              `json:"phoneNumber"`
	ProfileImage       string              `json:"profileImage"`
	AccountBalance     float64             `gorm:"default:0" json:"accountBalance"`
	RoleId             int                 `gorm:"not null;default:1" json:"roleId"`
	WatchList          []WatchListModel    `gorm:"foreignKey:userId" json:"watchList"`
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

func (u *User) BeforeUpdate(tx *gorm.DB) error {
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

func GetAllUsers(limit, offset int) ([]User, error) {
	var user []User
	if err := database.DB.
		Preload("Address").
		Preload("PortFolio").
		Preload("Transaction").
		Preload("Notification").
		Limit(limit).
		Offset(offset).
		Find(&user).Error; err != nil {
		log.Error().Err(err).Msg("issue persist in user_model/GetAllUsers")
		return nil, err
	}
	return user, nil

}

func GetUserById(id string) (*User, error) {
	var user User
	if err := database.DB.
		Where("id = ?", id).
		Preload("Address").
		Preload("PortFolio").
		Preload("Transaction").
		Preload("Notification").
		First(&user).Error; err != nil {
		log.Error().Err(err).Msg("issue persist in user_model/GetAllUsers")
		return nil, err
	}

	return &user, nil
}

func GetUserByEmail(email string) (*User, error) {
	var user User
	if err := database.DB.
		Where("email = ?", email).
		Preload("Address").
		Preload("PortFolio").
		Preload("Transaction").
		Preload("Notification").
		First(&user).Error; err != nil {
		log.Error().Err(err).Msg("issue persist in user_model/GetAllUsers")
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

func UpdateUserVerificationStatus(email string, status bool) error {
	if err := database.DB.Model(&User{}).
		Where("email = ?", email).
		Update("verification_status", status).Error; err != nil {
		log.Error().Err(err).Msg("issue lie in the user_model/UpdateUserVerificationStatus")
		return err
	}
	return nil
}

func UpdateLastLogin(email string, time time.Time) error {

	if err := database.DB.Where("email = ?", email).Update("last_login = ?", time).Error; err != nil {
		log.Error().Err(err).Msg("issue lie in the user_model/UpdateLastLogin")
		return err
	}
	return nil
}

func GetLastLogin(email string) (time.Time, error) {
	userByEmail, err := GetUserByEmail(email)
	if err != nil {
		log.Error().Err(err).Msg("issue lie in the user_model/GetLastLogin")
		return time.Time{}, err
	}
	return userByEmail.LastLogin, nil
}

func DeleteUserById(id string) error {
	return database.DB.Where("id = ?", id).Delete(&User{}).Error
}

func DeleteUserByEmail(email string) error {
	return database.DB.Where("email = ?", email).Delete(&User{}).Error
}

func UpdateUserLastLogin(email string) error {
	if err := database.DB.Model(&User{}).Where("email = ?", email).Update("last_login = ?", time.Now()).Error; err != nil {
		log.Error().Err(err).Msg("issue persist in the portfolio_model/UpdateUserLastLogin")
		return err
	}
	return nil
}
