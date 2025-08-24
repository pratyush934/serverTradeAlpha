package models

import "time"

type User struct {
	Id        string      `gorm:"primaryKey;type:varchar(151)" json:"id"`
	Name      string      `gorm:"not null" json:"name"`
	Email     string      `gorm:"not null;unique" json:"email"`
	PortFolio []ProtFolio `gorm:"foreignKey:userId" json:"portFolio"`
	RoleId    int         `gorm:"not null;default:1" json:"roleId"`
	Role      Role        `gorm:"not null;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"role"`
	IsActive  bool        `json:"isActive"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
}
