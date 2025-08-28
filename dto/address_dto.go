package dto

type AddressModelDTO struct {
	StreetName string ` json:"streetName"`
	LandMark   string ` json:"landMark"`
	ZipCode    string ` json:"zipCode"`
	City       string ` json:"city"`
	State      string ` json:"state"`
	Country    string `json:"country"`
}
