package repository

import (
	"context"
	"database/sql"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestRepository_GetUsersByPhoneNumberSuccess(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Fatal("error init mock", err)
	}
	defer db.Close()

	repo := Repository{
		Db: db,
	}

	mock.ExpectQuery("SELECT * FROM users WHERE phone_number = $1").WithArgs("628123456789").WillReturnRows(sqlmock.NewRows([]string{"id", "full_name", "phone_number", "password"}).AddRow(1, "Kabuki", "628123456789", "$2a$04$51FdsMsF1NCHjDP0VjApzO6o0Z.1Baf1nD5ua7CTO2Pmb0rTi2gs6"))

	_, err = repo.GetUsersByPhoneNumber(context.TODO(), "628123456789")
	if err != nil {
		t.Error("[TestRepository_GetUsersByPhoneNumberSuccess] assertion err failed. Expect nil, actual value: ", err)
	}
}

func TestRepository_GetUsersByPhoneNumberFailed(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Fatal("error init mock", err)
	}
	defer db.Close()

	repo := Repository{
		Db: db,
	}

	mock.ExpectQuery("SELECT * FROM users WHERE phone_number = $1").WithArgs("628123456789").WillReturnError(sql.ErrConnDone)

	_, err = repo.GetUsersByPhoneNumber(context.TODO(), "628123456789")
	if err == nil {
		t.Error("[TestRepository_GetUsersByPhoneNumberFailed] assertion err failed. Expect error is not nil")
	}
}

func TestRepository_CheckIfPhoneNumberExistFailedGetData(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Fatal("error init mock", err)
	}
	defer db.Close()

	repo := Repository{
		Db: db,
	}

	mock.ExpectQuery("SELECT id FROM users WHERE phone_number = $1").WithArgs("628123456789").WillReturnError(sql.ErrConnDone)

	_, err = repo.CheckIfPhoneNumberExist(context.TODO(), "628123456789")
	if err == nil {
		t.Error("[TestRepository_CheckIfPhoneNumberExistFailedGetData] assertion err failed. Expect error not nil")
	}
}

func TestRepository_CheckIfPhoneNumberExistSuccessIsNotExist(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Fatal("error init mock", err)
	}
	defer db.Close()

	repo := Repository{
		Db: db,
	}

	mock.ExpectQuery("SELECT id FROM users WHERE phone_number = $1").WithArgs("628123456789").WillReturnError(sql.ErrNoRows)

	_, err = repo.CheckIfPhoneNumberExist(context.TODO(), "628123456789")
	if err != nil {
		t.Error("[TestRepository_CheckIfPhoneNumberExistSuccessIsNotExist] assertion err failed. Expect error nil, actual error: ", err)
	}
}

func TestRepository_CheckIfPhoneNumberExistSuccessIsExist(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Fatal("error init mock", err)
	}
	defer db.Close()

	repo := Repository{
		Db: db,
	}

	mock.ExpectQuery("SELECT id FROM users WHERE phone_number = $1").WithArgs("628123456789").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	_, err = repo.CheckIfPhoneNumberExist(context.TODO(), "628123456789")
	if err != nil {
		t.Error("[TestRepository_CheckIfPhoneNumberExistSuccessIsExist] assertion err failed. Expect error nil, actual error: ", err)
	}
}

func TestRepository_RegisterFailed(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Fatal("error init mock", err)
	}
	defer db.Close()

	repo := Repository{
		Db: db,
	}

	dataTest := RegisterInput{
		FullName:    "Warga Konoha",
		PhoneNumber: "+628678990234",
		Password:    "$2a$04$lo.0ExyOiqlFFarEuQz35uExyLLaKZcPECT7zLbtgFzHyCNs0KP8.",
	}

	mock.ExpectQuery("INSERT INTO users(full_name, phone_number, password) VALUES ($1, $2, $3) RETURNING id").WithArgs(dataTest.FullName, dataTest.PhoneNumber, dataTest.Password).WillReturnError(sql.ErrConnDone)

	_, err = repo.Register(context.TODO(), dataTest)
	if err == nil {
		t.Error("[TestRepository_RegisterFailed] assertion err failed. Expect error not nil")
	}
}

func TestRepository_RegisterSuccess(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Fatal("error init mock", err)
	}
	defer db.Close()

	repo := Repository{
		Db: db,
	}

	dataTest := RegisterInput{
		FullName:    "Warga Konoha",
		PhoneNumber: "+628678990234",
		Password:    "$2a$04$lo.0ExyOiqlFFarEuQz35uExyLLaKZcPECT7zLbtgFzHyCNs0KP8.",
	}

	mock.ExpectQuery("INSERT INTO users(full_name, phone_number, password) VALUES ($1, $2, $3) RETURNING id").WithArgs(dataTest.FullName, dataTest.PhoneNumber, dataTest.Password).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	_, err = repo.Register(context.TODO(), dataTest)
	if err != nil {
		t.Error("[TestRepository_RegisterSuccess] assertion err failed. Expect error nil, actual err: ", err)
	}
}

func TestRepository_UpdateLoginCounterErrorGetData(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Fatal("error init mock", err)
	}
	defer db.Close()

	repo := Repository{
		Db: db,
	}

	dataTest := LoginHistoriesModel{
		ID:      1,
		UsersID: 1,
		Counter: 9,
	}

	mock.ExpectQuery("SELECT * FROM login_histories WHERE users_id = $1").WithArgs(dataTest.UsersID).WillReturnError(sql.ErrConnDone)

	err = repo.UpdateLoginCounter(context.TODO(), dataTest.UsersID)
	if err == nil {
		t.Error("[TestRepository_UpdateLoginCounterErrorGetData] assertion err failed. Expect error not nil")
	}
}

func TestRepository_UpdateLoginCounterErrorInsertNewData(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Fatal("error init mock", err)
	}
	defer db.Close()

	repo := Repository{
		Db: db,
	}

	dataTest := LoginHistoriesModel{
		ID:      1,
		UsersID: 1,
		Counter: 9,
	}

	mock.ExpectQuery("SELECT * FROM login_histories WHERE users_id = $1").WithArgs(dataTest.UsersID).WillReturnError(sql.ErrNoRows)
	mock.ExpectQuery("INSERT INTO login_histories(users_id, counter) VALUES($1, 1)").WithArgs(dataTest.UsersID).WillReturnError(sql.ErrNoRows)

	err = repo.UpdateLoginCounter(context.TODO(), dataTest.UsersID)
	if err == nil {
		t.Error("[TestRepository_UpdateLoginCounterErrorInsertNewData] assertion err failed. Expect error not nil")
	}
}

func TestRepository_UpdateLoginCounterSuccessInsertNewData(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Fatal("error init mock", err)
	}
	defer db.Close()

	repo := Repository{
		Db: db,
	}

	dataTest := LoginHistoriesModel{
		ID:      1,
		UsersID: 1,
		Counter: 9,
	}

	mock.ExpectQuery("SELECT * FROM login_histories WHERE users_id = $1").WithArgs(dataTest.UsersID).WillReturnError(sql.ErrNoRows)
	mock.ExpectQuery("INSERT INTO login_histories(users_id, counter) VALUES($1, 1)").WithArgs(dataTest.UsersID).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(dataTest.ID))

	err = repo.UpdateLoginCounter(context.TODO(), dataTest.UsersID)
	if err != nil {
		t.Error("[TestRepository_UpdateLoginCounterSuccessInsertNewData] assertion err failed. Expect error nil, actual: ", err)
	}
}

func TestRepository_UpdateLoginCounterErrorUpdateData(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Fatal("error init mock", err)
	}
	defer db.Close()

	repo := Repository{
		Db: db,
	}

	dataTest := LoginHistoriesModel{
		ID:      1,
		UsersID: 1,
		Counter: 9,
	}

	mock.ExpectQuery("SELECT * FROM login_histories WHERE users_id = $1").WithArgs(dataTest.UsersID).WillReturnRows(sqlmock.NewRows([]string{"id", "users_id", "counter"}).AddRow(dataTest.ID, dataTest.UsersID, dataTest.Counter))
	mock.ExpectQuery("UPDATE login_histories SET counter = $1 WHERE users_id = $2").WithArgs(dataTest.Counter, dataTest.UsersID).WillReturnError(sql.ErrConnDone)

	err = repo.UpdateLoginCounter(context.TODO(), dataTest.UsersID)
	if err == nil {
		t.Error("[TestRepository_UpdateLoginCounterErrorUpdateData] assertion err failed. Expect error not nil")
	}
}

func TestRepository_UpdateLoginCounterSuccessUpdateData(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Fatal("error init mock", err)
	}
	defer db.Close()

	repo := Repository{
		Db: db,
	}

	dataTest := LoginHistoriesModel{
		ID:      1,
		UsersID: 1,
		Counter: 9,
	}

	mock.ExpectQuery("SELECT * FROM login_histories WHERE users_id = $1").WithArgs(dataTest.UsersID).WillReturnRows(sqlmock.NewRows([]string{"id", "users_id", "counter"}).AddRow(dataTest.ID, dataTest.UsersID, dataTest.Counter))
	mock.ExpectQuery("UPDATE login_histories SET counter = $1 WHERE users_id = $2").WithArgs(dataTest.Counter+1, dataTest.UsersID).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(dataTest.ID))

	err = repo.UpdateLoginCounter(context.TODO(), dataTest.UsersID)
	if err != nil {
		t.Error("[TestRepository_UpdateLoginCounterSuccessUpdateData] assertion err failed. Expect error nil, actual: ", err)
	}
}

func TestRepository_UpdateProfileSuccessFullNameOnly(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Fatal("error init mock", err)
	}
	defer db.Close()

	repo := Repository{
		Db: db,
	}

	dataTest := UpdateProfileInput{
		FullName: "Kabuki",
		// PhoneNumber: "+628678990234",
	}

	mock.ExpectQuery("UPDATE users SET full_name = Kabuki").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	_, err = repo.UpdateProfile(context.TODO(), dataTest)
	if err != nil {
		t.Error("[TestRepository_UpdateProfileSuccessFullNameOnly] assertion err failed. Expect error nil, actual: ", err)
	}
}

func TestRepository_UpdateProfileSuccessPhoneNumberOnly(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Fatal("error init mock", err)
	}
	defer db.Close()

	repo := Repository{
		Db: db,
	}

	dataTest := UpdateProfileInput{
		// FullName: "Kabuki",
		PhoneNumber: "+628678990234",
	}

	mock.ExpectQuery("UPDATE users SET phone_number = +628678990234").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	_, err = repo.UpdateProfile(context.TODO(), dataTest)
	if err != nil {
		t.Error("[TestRepository_UpdateProfileSuccessPhoneNumberOnly] assertion err failed. Expect error nil, actual: ", err)
	}
}

func TestRepository_UpdateProfileSuccessFullNameAndPhoneNumber(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Fatal("error init mock", err)
	}
	defer db.Close()

	repo := Repository{
		Db: db,
	}

	dataTest := UpdateProfileInput{
		FullName:    "Kabuki",
		PhoneNumber: "+628678990234",
	}

	mock.ExpectQuery("UPDATE users SET full_name = Kabuki, phone_number = +628678990234").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	_, err = repo.UpdateProfile(context.TODO(), dataTest)
	if err != nil {
		t.Error("[TestRepository_UpdateProfileSuccessFullNameAndPhoneNumber] assertion err failed. Expect error nil, actual: ", err)
	}
}

func TestRepository_UpdateProfileError(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Fatal("error init mock", err)
	}
	defer db.Close()

	repo := Repository{
		Db: db,
	}

	dataTest := UpdateProfileInput{
		FullName:    "Kabuki",
		PhoneNumber: "+628678990234",
	}

	mock.ExpectQuery("UPDATE users SET full_name = Kabuki, phone_number = +628678990234").WillReturnError(sql.ErrConnDone)

	_, err = repo.UpdateProfile(context.TODO(), dataTest)
	if err == nil {
		t.Error("[TestRepository_UpdateProfileError] assertion err failed. Expect error not nil")
	}
}
