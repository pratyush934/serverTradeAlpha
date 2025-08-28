package dto

type NotificationDTO struct {
	Message    string `json:"message"`
	ReadStatus bool   `json:"readStatus"`
}
