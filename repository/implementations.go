package repository

import (
	"context"
)

func (r *Repository) GetTestById(ctx context.Context, input GetTestByIdInput) (output GetTestByIdOutput, err error) {
	err = r.Db.QueryRowContext(ctx, "SELECT name FROM test WHERE id = $1", input.Id).Scan(&output.Name)
	if err != nil {
		return
	}
	return
}

func (r *Repository) Register(ctx context.Context, input RegisterInput) (output RegisterOutput, err error) {
	return
}

func (r *Repository) GetUsersByPhoneNumber(ctx context.Context, input string) (output GetUsersByPhoneNumberOutput, err error) {
	err = r.Db.QueryRowContext(ctx, "SELECT * FROM users WHERE phone_number = $1", input).Scan(&output.ID, &output.FullName, &output.PhoneNumber, &output.Password, &output.CreatedAt)
	if err != nil {
		return
	}

	return
}
