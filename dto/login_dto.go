package dto

type LoginModel struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	OAuthId  string `json:"oauth_id"`
	Provider string `json:"provider"`
	Image    string `json:"image"`
}
