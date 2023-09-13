// This file contains types that are used in the repository layer.
package repository

import "time"

type GetTestByIdInput struct {
	Id string
}

type GetTestByIdOutput struct {
	Name string
}

type RegisterInput struct {
	FullName    string `json:"full_name"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type RegisterOutput struct {
	ID int `json:"id"`
}

type GetUsersByPhoneNumberOutput struct {
	ID          int       `json:"id"`
	FullName    string    `json:"full_name"`
	PhoneNumber string    `json:"phone_number"`
	Password    string    `json:"password"`
	CreatedAt   time.Time `json:"created_at"`
}

type LoginHistoriesModel struct {
	ID      int `json:"id"`
	UsersID int `json:"users_id"`
	Counter int `json:"counter"`
}

type UpdateProfileInput struct {
	FullName    string
	PhoneNumber string
}

type UpdateProfileOutput struct {
	Success bool
}
