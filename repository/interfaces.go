// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import (
	"context"
)

type RepositoryInterface interface {
	Register(ctx context.Context, input RegisterInput) (output RegisterOutput, err error)
	GetUsersByPhoneNumber(ctx context.Context, userPhoneNumber string) (output GetUsersByPhoneNumberOutput, err error)
	UpdateLoginCounter(ctx context.Context, input int) error
	CheckIfPhoneNumberExist(ctx context.Context, input string) (isExist bool, err error)
	UpdateProfile(ctx context.Context, input UpdateProfileInput) (success bool, err error)
}
