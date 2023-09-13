package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/triardn/panganggo/commons"
	"github.com/triardn/panganggo/generated"
	"github.com/triardn/panganggo/repository"
)

// This is just a test endpoint to get you started. Please delete this endpoint.
// (GET /hello)
func (s *Server) Hello(ctx echo.Context, params generated.HelloParams) error {
	var resp generated.HelloResponse
	resp.Message = fmt.Sprintf("Hello User %d", params.Id)
	return ctx.JSON(http.StatusOK, resp)
}

func (s *Server) Registration(ctx echo.Context) error {
	registerRequest := generated.RegisterRequest{}
	if err := ctx.Bind(&registerRequest); err != nil {
		return err
	}

	// validate request
	errors := validateRequest(registerRequest)
	if len(errors) > 0 {
		resp := generated.BadRequestResponse{
			Message: "Gagal melakukan registrasi. Periksa detail berikut:",
			Detail:  errors,
		}

		return ctx.JSON(http.StatusBadRequest, resp)
	}

	// create new user
	payload := repository.RegisterInput{
		FullName:    registerRequest.FullName,
		PhoneNumber: registerRequest.PhoneNumber,
		Password:    commons.HashAndSalt([]byte(registerRequest.Password)),
	}

	data, err := s.Repository.Register(ctx.Request().Context(), payload)
	if err != nil {
		// TODO: add log
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: "Terjadi kesalahan pada server. Silakan coba lagi."})
	}

	return ctx.JSON(http.StatusOK, data)
}

func (s *Server) Login(ctx echo.Context) error {
	loginRequest := generated.LoginRequest{}
	if err := ctx.Bind(&loginRequest); err != nil {
		return err
	}

	// get user by phone_number
	users, err := s.Repository.GetUsersByPhoneNumber(ctx.Request().Context(), loginRequest.PhoneNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: "User tidak ditemukan."})
		}

		// TODO: add log
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: "Gagal mendapatkan informasi user dari server. Silakan coba lagi."})
	}

	// compare password
	if !commons.ComparePasswords(users.Password, []byte(loginRequest.Password)) {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: "Phone number atau password salah."})
	}

	// compose JWT token
	token, err := s.JWT.CreateToken(ctx.Request().Context(), commons.UserData{ID: users.ID, FullName: users.FullName, PhoneNumber: users.PhoneNumber})
	if err != nil {
		// TODO: add log
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: "Terjadi kesalahan pada sistem. Silakan coba lagi."})
	}

	// increase login counter
	err = s.Repository.UpdateLoginCounter(ctx.Request().Context(), int(users.ID))
	if err != nil {
		// TODO: add log
		fmt.Println(err)
	}

	return ctx.JSON(http.StatusOK, generated.LoginResponse{Id: int(users.ID), Token: token})
}

func (s *Server) GetUserDetailByID(ctx echo.Context, id int) error {
	ok, data := validateToken(ctx, s.JWT)
	if !ok {
		return ctx.JSON(http.StatusForbidden, generated.ErrorResponse{Message: "Forbidden."})
	}

	return ctx.JSON(http.StatusOK, generated.UserDetailResponse{FullName: data.FullName, PhoneNumber: data.PhoneNumber})
}

func (s *Server) UpdateProfile(ctx echo.Context) error {
	ok, _ := validateToken(ctx, s.JWT)
	if !ok {
		return ctx.JSON(http.StatusForbidden, generated.ErrorResponse{Message: "Forbidden."})
	}

	updateProfileRequest := generated.UpdateProfileRequest{}
	if err := ctx.Bind(&updateProfileRequest); err != nil {
		return err
	}

	// validate phone number

	// check if phone number is exist
	isExist, err := s.Repository.CheckIfPhoneNumberExist(ctx.Request().Context(), *updateProfileRequest.PhoneNumber)
	if err != nil {
		// TODO: add log
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: "Terjadi kesalahan pada sistem. Silakan coba lagi."})
	}

	if isExist {
		return ctx.JSON(http.StatusConflict, generated.ErrorResponse{Message: "Conflict."})
	}

	// update profile
	payload := repository.UpdateProfileInput{}

	if updateProfileRequest.FullName != nil {
		payload.FullName = *updateProfileRequest.FullName
	}

	if updateProfileRequest.PhoneNumber != nil {
		payload.PhoneNumber = *updateProfileRequest.PhoneNumber
	}

	_, err = s.Repository.UpdateProfile(ctx.Request().Context(), payload)
	if err != nil {
		// TODO: add log
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: "Gagal update profile. Silakan coba lagi."})
	}

	return ctx.JSON(http.StatusOK, generated.UpdateProfileResponse{Message: "Sukses update profile."})
}

func validateToken(ctx echo.Context, jwt commons.JWT) (isValid bool, data commons.UserData) {
	// get token
	authBearer := ctx.Request().Header.Get("Authorization")
	temp := strings.Split(authBearer, " ")

	tempData, err := jwt.ValidateToken(temp[1])
	if err != nil {
		// TODO: log error
		fmt.Println(err)
		return
	}

	// map data
	data.ID = tempData.(commons.UserData).ID
	data.FullName = tempData.(commons.UserData).FullName
	data.PhoneNumber = tempData.(commons.UserData).PhoneNumber

	isValid = true

	return
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func validateFullName(fullName string) (isValid bool, errors ValidationError) {
	if len(fullName) < 3 || len(fullName) > 60 {
		errors.Field = "full_name"
		errors.Message = "karakter fullname minimal 3 karakter dan maksimal 60 karakter."

		return
	}

	return true, errors
}

func validatePassword(password string) (isValid bool, errors ValidationError) {
	isValid = commons.ValidatePassword(password)
	if !isValid {
		errors.Field = "password"
		errors.Message = "password tidak sesuai dengan ketentuan"

		return
	}

	return true, errors
}

func validatePhoneNumber(phoneNumber string) (isValid bool, errors ValidationError) {
	var message []string
	if len(phoneNumber) < 10 || len(phoneNumber) > 13 {
		message = append(message, "nomor telpon minimal 10 karakter dan maksimal 13 karakter")
	}

	if len(phoneNumber) >= 10 {
		prefix := phoneNumber[0:3]
		if prefix != "+62" {
			message = append(message, "nomor telpon harus diawali dengan +62")
		}
	}

	// any validation errors
	if len(message) > 0 {
		errors.Field = "phone_number"
		if len(message) > 1 {
			errors.Message = strings.Join(message, " dan")
		} else {
			errors.Message = message[0]
		}

		return
	}

	return true, errors
}

func validateRequest(input generated.RegisterRequest) (errors []string) {
	var validationErrors []ValidationError

	// validate fullname
	isValid, error := validateFullName(input.FullName)
	if !isValid {
		validationErrors = append(validationErrors, error)
	}

	// validate phone number
	isValid, error = validatePhoneNumber(input.PhoneNumber)
	if !isValid {
		validationErrors = append(validationErrors, error)
	} else {
		// check if phone number is already used
	}

	// validate password
	isValid, error = validatePassword(input.Password)
	if !isValid {
		validationErrors = append(validationErrors, error)
	}

	// iterate all errors
	if len(validationErrors) > 0 {
		for _, data := range validationErrors {
			errors = append(errors, fmt.Sprintf("%s: %s", data.Field, data.Message))
		}
	}

	return
}
