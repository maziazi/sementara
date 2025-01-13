package model

import "time"

type Managers struct {
	ID              uint      `json:"id"`
	Email           string    `json:"email"`
	Password        string    `json:"-"` // Jangan kembalikan password di response
	CreatedAt       time.Time `json:"created_at"`
	Name            string    `json:"name"`
	CompanyName     string    `json:"company_name"`
	ManagerImageURI string    `json:"manager_image_uri"`
	CompanyImageURI string    `json:"company_image_uri"`
}
