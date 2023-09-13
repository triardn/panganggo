package repository

import (
	"context"
	"database/sql"
	"fmt"
)

func (r *Repository) GetTestById(ctx context.Context, input GetTestByIdInput) (output GetTestByIdOutput, err error) {
	err = r.Db.QueryRowContext(ctx, "SELECT name FROM test WHERE id = $1", input.Id).Scan(&output.Name)
	if err != nil {
		return
	}
	return
}

func (r *Repository) Register(ctx context.Context, input RegisterInput) (output RegisterOutput, err error) {
	sqlStatement := `
		INSERT INTO users(full_name, phone_number, password)
		VALUES ($1, $2, $3)
		RETURNING id`
	err = r.Db.QueryRowContext(ctx, sqlStatement, input.FullName, input.PhoneNumber, input.Password).Scan(&output.ID)
	if err != nil {
		return
	}

	return
}

func (r *Repository) GetUsersByPhoneNumber(ctx context.Context, input string) (output GetUsersByPhoneNumberOutput, err error) {
	err = r.Db.QueryRowContext(ctx, "SELECT * FROM users WHERE phone_number = $1", input).Scan(&output.ID, &output.FullName, &output.PhoneNumber, &output.Password, &output.CreatedAt)
	if err != nil {
		return
	}

	return
}

func (r *Repository) UpdateLoginCounter(ctx context.Context, input int) error {
	model := LoginHistoriesModel{}
	err := r.Db.QueryRowContext(ctx, "SELECT * FROM login_histories WHERE users_id = $1", input).Scan(&model.ID, &model.UsersID, &model.Counter)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if err == sql.ErrNoRows {
		// create new record
		err = r.Db.QueryRowContext(ctx, "INSERT INTO login_histories(users_id, counter) VALUES($1, 1)", input).Err()
		if err != nil {
			return err
		}

		return nil
	}

	// update directly
	err = r.Db.QueryRowContext(ctx, "UPDATE login_histories SET counter = $1 WHERE users_id = $2", model.Counter+1, input).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) CheckIfPhoneNumberExist(ctx context.Context, input string) (isExist bool, err error) {
	var tempID int
	err = r.Db.QueryRowContext(ctx, "SELECT id FROM users WHERE phone_number = $1", input).Scan(&tempID)
	if err != nil && err != sql.ErrNoRows {
		return
	}

	if err == sql.ErrNoRows {
		return false, nil
	}

	return true, nil
}

func (r *Repository) UpdateProfile(ctx context.Context, input UpdateProfileInput) (success bool, err error) {
	query := "UPDATE users SET"

	if input.FullName != "" {
		query = query + fmt.Sprintf(" full_name = %s", input.FullName)
	}

	if input.PhoneNumber != "" {
		if input.FullName != "" {
			query = query + fmt.Sprintf(", phone_number = %s", input.PhoneNumber)
		} else {
			query = query + fmt.Sprintf(" phone_number = %s", input.PhoneNumber)
		}
	}

	fmt.Println(query)

	err = r.Db.QueryRowContext(ctx, query).Err()
	if err != nil {
		return
	}

	return true, nil
}
