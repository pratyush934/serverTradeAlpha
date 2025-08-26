package models

import "time"

type Role struct {
	Id          int    `json:"id"`
	RoleName    string `json:"roleName"`
	Description string `json:"description"`
	CreatedAt   time.Time
}

func CreateRole() {

}
