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
	return nil
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
