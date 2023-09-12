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
}

type RegisterOutput struct {
}

type GetUsersByPhoneNumberOutput struct {
	ID          uint64    `json:"id"`
	FullName    string    `json:"full_name"`
	PhoneNumber string    `json:"phone_number"`
	Password    string    `json:"password"`
	CreatedAt   time.Time `json:"created_at"`
}
