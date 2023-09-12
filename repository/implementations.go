package repository

import (
	"context"
	"database/sql"
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
