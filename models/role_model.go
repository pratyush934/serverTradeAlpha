package models

type Role struct {
	Id          int    `json:"id"`
	RoleName    string `json:"roleName"`
	Description string `json:"description"`
}
