package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/triardn/panganggo/commons"
	"github.com/triardn/panganggo/repository"
	"go.uber.org/mock/gomock"
)

func TestHello(t *testing.T) {

}

var (
	registrationJSON = `{"full_name": "Petrikor", "phone_number": "+6281234567", "password": "Petr1k0r@10ss"}`
)

func TestRegister_FailedAlreadyResgitered(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/registration", strings.NewReader(registrationJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockCtrl := gomock.NewController(t)
	mockRepository := repository.NewMockRepositoryInterface(mockCtrl)

	mockJWT := commons.NewMockJWT(mockCtrl)

	mockRepository.EXPECT().CheckIfPhoneNumberExist(c.Request().Context(), "+6281234567").Return(true, nil)

	// // create new user
	// payload := repository.RegisterInput{
	// 	FullName:    "Petrikor",
	// 	PhoneNumber: "+6281234567",
	// 	Password:    commons.HashAndSalt([]byte("Petr1k0r@10ss")),
	// }

	// mockRepository.EXPECT().Register(c.Request().Context(), payload).Return(repository.RegisterOutput{}, sql.ErrConnDone)

	server := NewServer(NewServerOptions{
		Repository: mockRepository,
		JWT:        mockJWT,
	})

	err := server.Registration(c)

	assert.Nil(t, err)
}

func TestRegister_FailedRequestBodyNotValid(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/registration", strings.NewReader("1"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockCtrl := gomock.NewController(t)
	mockRepository := repository.NewMockRepositoryInterface(mockCtrl)

	mockJWT := commons.NewMockJWT(mockCtrl)

	server := NewServer(NewServerOptions{
		Repository: mockRepository,
		JWT:        mockJWT,
	})

	err := server.Registration(c)

	assert.NotNil(t, err)
}

func TestRegister_FailedPasswordDidNotMatchRule(t *testing.T) {
	registrationJSONFailed := `{"full_name": "Petrikor", "phone_number": "+6281234567", "password": "petr1k0r@10ss"}`

	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/registration", strings.NewReader(registrationJSONFailed))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockCtrl := gomock.NewController(t)
	mockRepository := repository.NewMockRepositoryInterface(mockCtrl)

	mockJWT := commons.NewMockJWT(mockCtrl)

	server := NewServer(NewServerOptions{
		Repository: mockRepository,
		JWT:        mockJWT,
	})

	err := server.Registration(c)

	assert.Nil(t, err)
}

func TestServer_LoginRequestBodyNotValid(t *testing.T) {
	// loginJson := `{"phone_number": "+6281234567", "password": "petr1k0r@10ss"}`

	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader("!"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockCtrl := gomock.NewController(t)
	mockRepository := repository.NewMockRepositoryInterface(mockCtrl)

	mockJWT := commons.NewMockJWT(mockCtrl)

	server := NewServer(NewServerOptions{
		Repository: mockRepository,
		JWT:        mockJWT,
	})

	err := server.Login(c)

	assert.NotNil(t, err)
}

func TestServer_LoginUserNotFound(t *testing.T) {
	loginJson := `{"phone_number": "+628678990234", "password": "$2a$04$lo.0ExyOiqlFFarEuQz35uExyLLaKZcPECT7zLbtgFzHyCNs0KP8."}`

	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(loginJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockCtrl := gomock.NewController(t)
	mockRepository := repository.NewMockRepositoryInterface(mockCtrl)

	mockJWT := commons.NewMockJWT(mockCtrl)

	mockRepository.EXPECT().GetUsersByPhoneNumber(c.Request().Context(), "+628678990234").Return(repository.GetUsersByPhoneNumberOutput{}, sql.ErrNoRows)

	server := NewServer(NewServerOptions{
		Repository: mockRepository,
		JWT:        mockJWT,
	})

	err := server.Login(c)

	assert.Nil(t, err)
}

func TestServer_LoginGetUserGotInternalError(t *testing.T) {
	loginJson := `{"phone_number": "+628678990234", "password": "$2a$04$lo.0ExyOiqlFFarEuQz35uExyLLaKZcPECT7zLbtgFzHyCNs0KP8."}`

	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(loginJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockCtrl := gomock.NewController(t)
	mockRepository := repository.NewMockRepositoryInterface(mockCtrl)

	mockJWT := commons.NewMockJWT(mockCtrl)

	mockRepository.EXPECT().GetUsersByPhoneNumber(c.Request().Context(), "+628678990234").Return(repository.GetUsersByPhoneNumberOutput{}, sql.ErrConnDone)

	server := NewServer(NewServerOptions{
		Repository: mockRepository,
		JWT:        mockJWT,
	})

	err := server.Login(c)

	assert.Nil(t, err)
}

func TestServer_GetUserDetailByIDError(t *testing.T) {
	token := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwiZnVsbF9uYW1lIjoiV2FyZ2EgS29ub2hhIiwicGhvbmVfbnVtYmVyIjoiNjI4MTIzNDU2Nzg5IiwiZXhwIjoxNjk0ODYxMzE4fQ.S7Ang3pIbMbU4NOPpb2360Qrxgs1kAE2Za20SOL_eC3hOaowpDfYtZ8JAaDFt-NdnU_S5D0yae-3s9n0GjwjmzBbcCzgQmNxRsP4HO1OZ0qJP688Wd1WQftuTapjblUrPPsgsaEiixUcDhLCDm5DE9vcUA-6BTYx3C7m3UeHxBhQX4oXqcFioMZg_qja3d2qSFlLjTEwnbuDL5mRaz2ac3cjB6Sv3xZf2aGGfrnPsnAKziX0wr3WN8H6QApi1SZLnYlrS89C5sIUyVSvVj4o5cU8u-2j6KXQfy3smYVBE4Hm5MGg4_395TjfIGdPf5fPc6IjPArCqpEkGlSMKIcXKIutW1Kv3WDhlAJ_WozN2aa6W4glNvcE3CSBqIqZbYzovncSEVEF-SePU8T8OI1zEjiwdmWX_jeCJMqUu-IdndHdtR9kf5N-m3SMui4N4F7hovjJNCh_0g1BjA0OaHJV62r1l7S1w0tY-vXkcglBLc_eRLFb6-zbhNGtKhKKDu2YgbfPoP5f8aivTWXNcXDV8NrXZuKXCuwhM1FU7lDj4FpYTerWzETXTZv3tzgSPAdO3coUwuhMlSVjMdb8bYzW5LHl6nQlATbj6FtNChCdaoymFp6nPm7li1IgqYluew1ahpoFCE8sFSbw2Z4__DHSSPiJR6JErFV3phuGk04CY-k"

	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", token))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockCtrl := gomock.NewController(t)
	mockRepository := repository.NewMockRepositoryInterface(mockCtrl)

	mockJWT := commons.NewMockJWT(mockCtrl)

	mockJWT.EXPECT().ValidateToken(token).Return(nil, errors.New("some-error"))

	server := NewServer(NewServerOptions{
		Repository: mockRepository,
		JWT:        mockJWT,
	})

	err := server.GetUserDetailByID(c, 1)

	assert.Nil(t, err)
}

func TestServer_GetUserDetailByIDSuccess(t *testing.T) {
	token := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwiZnVsbF9uYW1lIjoiV2FyZ2EgS29ub2hhIiwicGhvbmVfbnVtYmVyIjoiNjI4MTIzNDU2Nzg5IiwiZXhwIjoxNjk0ODYxMzE4fQ.S7Ang3pIbMbU4NOPpb2360Qrxgs1kAE2Za20SOL_eC3hOaowpDfYtZ8JAaDFt-NdnU_S5D0yae-3s9n0GjwjmzBbcCzgQmNxRsP4HO1OZ0qJP688Wd1WQftuTapjblUrPPsgsaEiixUcDhLCDm5DE9vcUA-6BTYx3C7m3UeHxBhQX4oXqcFioMZg_qja3d2qSFlLjTEwnbuDL5mRaz2ac3cjB6Sv3xZf2aGGfrnPsnAKziX0wr3WN8H6QApi1SZLnYlrS89C5sIUyVSvVj4o5cU8u-2j6KXQfy3smYVBE4Hm5MGg4_395TjfIGdPf5fPc6IjPArCqpEkGlSMKIcXKIutW1Kv3WDhlAJ_WozN2aa6W4glNvcE3CSBqIqZbYzovncSEVEF-SePU8T8OI1zEjiwdmWX_jeCJMqUu-IdndHdtR9kf5N-m3SMui4N4F7hovjJNCh_0g1BjA0OaHJV62r1l7S1w0tY-vXkcglBLc_eRLFb6-zbhNGtKhKKDu2YgbfPoP5f8aivTWXNcXDV8NrXZuKXCuwhM1FU7lDj4FpYTerWzETXTZv3tzgSPAdO3coUwuhMlSVjMdb8bYzW5LHl6nQlATbj6FtNChCdaoymFp6nPm7li1IgqYluew1ahpoFCE8sFSbw2Z4__DHSSPiJR6JErFV3phuGk04CY-k"

	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", token))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockCtrl := gomock.NewController(t)
	mockRepository := repository.NewMockRepositoryInterface(mockCtrl)

	mockJWT := commons.NewMockJWT(mockCtrl)

	mockJWT.EXPECT().ValidateToken(token).Return(commons.UserData{ID: 1, FullName: "Warga Konoha", PhoneNumber: "+628123456789"}, nil)

	server := NewServer(NewServerOptions{
		Repository: mockRepository,
		JWT:        mockJWT,
	})

	err := server.GetUserDetailByID(c, 1)

	assert.Nil(t, err)
}

func TestServer_UpdateProfileTokenNotValid(t *testing.T) {
	token := "1"

	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/users", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", token))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockCtrl := gomock.NewController(t)
	mockRepository := repository.NewMockRepositoryInterface(mockCtrl)

	mockJWT := commons.NewMockJWT(mockCtrl)

	mockJWT.EXPECT().ValidateToken(token).Return(nil, errors.New("some-error"))

	server := NewServer(NewServerOptions{
		Repository: mockRepository,
		JWT:        mockJWT,
	})

	err := server.UpdateProfile(c)

	assert.Nil(t, err)
}

func TestServer_UpdateProfileRequestBodyNotValid(t *testing.T) {
	token := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwiZnVsbF9uYW1lIjoiV2FyZ2EgS29ub2hhIiwicGhvbmVfbnVtYmVyIjoiNjI4MTIzNDU2Nzg5IiwiZXhwIjoxNjk0ODYxMzE4fQ.S7Ang3pIbMbU4NOPpb2360Qrxgs1kAE2Za20SOL_eC3hOaowpDfYtZ8JAaDFt-NdnU_S5D0yae-3s9n0GjwjmzBbcCzgQmNxRsP4HO1OZ0qJP688Wd1WQftuTapjblUrPPsgsaEiixUcDhLCDm5DE9vcUA-6BTYx3C7m3UeHxBhQX4oXqcFioMZg_qja3d2qSFlLjTEwnbuDL5mRaz2ac3cjB6Sv3xZf2aGGfrnPsnAKziX0wr3WN8H6QApi1SZLnYlrS89C5sIUyVSvVj4o5cU8u-2j6KXQfy3smYVBE4Hm5MGg4_395TjfIGdPf5fPc6IjPArCqpEkGlSMKIcXKIutW1Kv3WDhlAJ_WozN2aa6W4glNvcE3CSBqIqZbYzovncSEVEF-SePU8T8OI1zEjiwdmWX_jeCJMqUu-IdndHdtR9kf5N-m3SMui4N4F7hovjJNCh_0g1BjA0OaHJV62r1l7S1w0tY-vXkcglBLc_eRLFb6-zbhNGtKhKKDu2YgbfPoP5f8aivTWXNcXDV8NrXZuKXCuwhM1FU7lDj4FpYTerWzETXTZv3tzgSPAdO3coUwuhMlSVjMdb8bYzW5LHl6nQlATbj6FtNChCdaoymFp6nPm7li1IgqYluew1ahpoFCE8sFSbw2Z4__DHSSPiJR6JErFV3phuGk04CY-k"

	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/users", strings.NewReader("q"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", token))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockCtrl := gomock.NewController(t)
	mockRepository := repository.NewMockRepositoryInterface(mockCtrl)

	mockJWT := commons.NewMockJWT(mockCtrl)

	mockJWT.EXPECT().ValidateToken(token).Return(commons.UserData{ID: 1, FullName: "Warga Konoha", PhoneNumber: "+628123456789"}, nil)

	server := NewServer(NewServerOptions{
		Repository: mockRepository,
		JWT:        mockJWT,
	})

	err := server.UpdateProfile(c)

	assert.NotNil(t, err)
}

func TestServer_UpdateProfileGotInternalError(t *testing.T) {
	token := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwiZnVsbF9uYW1lIjoiV2FyZ2EgS29ub2hhIiwicGhvbmVfbnVtYmVyIjoiNjI4MTIzNDU2Nzg5IiwiZXhwIjoxNjk0ODYxMzE4fQ.S7Ang3pIbMbU4NOPpb2360Qrxgs1kAE2Za20SOL_eC3hOaowpDfYtZ8JAaDFt-NdnU_S5D0yae-3s9n0GjwjmzBbcCzgQmNxRsP4HO1OZ0qJP688Wd1WQftuTapjblUrPPsgsaEiixUcDhLCDm5DE9vcUA-6BTYx3C7m3UeHxBhQX4oXqcFioMZg_qja3d2qSFlLjTEwnbuDL5mRaz2ac3cjB6Sv3xZf2aGGfrnPsnAKziX0wr3WN8H6QApi1SZLnYlrS89C5sIUyVSvVj4o5cU8u-2j6KXQfy3smYVBE4Hm5MGg4_395TjfIGdPf5fPc6IjPArCqpEkGlSMKIcXKIutW1Kv3WDhlAJ_WozN2aa6W4glNvcE3CSBqIqZbYzovncSEVEF-SePU8T8OI1zEjiwdmWX_jeCJMqUu-IdndHdtR9kf5N-m3SMui4N4F7hovjJNCh_0g1BjA0OaHJV62r1l7S1w0tY-vXkcglBLc_eRLFb6-zbhNGtKhKKDu2YgbfPoP5f8aivTWXNcXDV8NrXZuKXCuwhM1FU7lDj4FpYTerWzETXTZv3tzgSPAdO3coUwuhMlSVjMdb8bYzW5LHl6nQlATbj6FtNChCdaoymFp6nPm7li1IgqYluew1ahpoFCE8sFSbw2Z4__DHSSPiJR6JErFV3phuGk04CY-k"

	updateProfileJSON := `{"full_name": "Kodok", "phone_number": "+628123456709"}`

	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/users", strings.NewReader(updateProfileJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", token))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockCtrl := gomock.NewController(t)
	mockRepository := repository.NewMockRepositoryInterface(mockCtrl)

	mockJWT := commons.NewMockJWT(mockCtrl)

	mockJWT.EXPECT().ValidateToken(token).Return(commons.UserData{ID: 1, FullName: "Warga Konoha", PhoneNumber: "+628123456789"}, nil)

	mockRepository.EXPECT().CheckIfPhoneNumberExist(c.Request().Context(), "+628123456709").Return(false, nil)

	// fullname := "Kodok"
	// phoneNumber := "+628123456709"
	// payload := generated.UpdateProfileRequest{
	// 	FullName:    &fullname,
	// 	PhoneNumber: &phoneNumber,
	// }

	payload := repository.UpdateProfileInput{
		FullName:    "Kodok",
		PhoneNumber: "+628123456709",
	}

	mockRepository.EXPECT().UpdateProfile(c.Request().Context(), payload).Return(false, sql.ErrConnDone)

	server := NewServer(NewServerOptions{
		Repository: mockRepository,
		JWT:        mockJWT,
	})

	err := server.UpdateProfile(c)

	assert.Nil(t, err)
}

func TestServer_UpdateProfileSuccess(t *testing.T) {
	token := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwiZnVsbF9uYW1lIjoiV2FyZ2EgS29ub2hhIiwicGhvbmVfbnVtYmVyIjoiNjI4MTIzNDU2Nzg5IiwiZXhwIjoxNjk0ODYxMzE4fQ.S7Ang3pIbMbU4NOPpb2360Qrxgs1kAE2Za20SOL_eC3hOaowpDfYtZ8JAaDFt-NdnU_S5D0yae-3s9n0GjwjmzBbcCzgQmNxRsP4HO1OZ0qJP688Wd1WQftuTapjblUrPPsgsaEiixUcDhLCDm5DE9vcUA-6BTYx3C7m3UeHxBhQX4oXqcFioMZg_qja3d2qSFlLjTEwnbuDL5mRaz2ac3cjB6Sv3xZf2aGGfrnPsnAKziX0wr3WN8H6QApi1SZLnYlrS89C5sIUyVSvVj4o5cU8u-2j6KXQfy3smYVBE4Hm5MGg4_395TjfIGdPf5fPc6IjPArCqpEkGlSMKIcXKIutW1Kv3WDhlAJ_WozN2aa6W4glNvcE3CSBqIqZbYzovncSEVEF-SePU8T8OI1zEjiwdmWX_jeCJMqUu-IdndHdtR9kf5N-m3SMui4N4F7hovjJNCh_0g1BjA0OaHJV62r1l7S1w0tY-vXkcglBLc_eRLFb6-zbhNGtKhKKDu2YgbfPoP5f8aivTWXNcXDV8NrXZuKXCuwhM1FU7lDj4FpYTerWzETXTZv3tzgSPAdO3coUwuhMlSVjMdb8bYzW5LHl6nQlATbj6FtNChCdaoymFp6nPm7li1IgqYluew1ahpoFCE8sFSbw2Z4__DHSSPiJR6JErFV3phuGk04CY-k"

	updateProfileJSON := `{"full_name": "Kodok", "phone_number": "+628123456709"}`

	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/users", strings.NewReader(updateProfileJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", token))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockCtrl := gomock.NewController(t)
	mockRepository := repository.NewMockRepositoryInterface(mockCtrl)

	mockJWT := commons.NewMockJWT(mockCtrl)

	mockJWT.EXPECT().ValidateToken(token).Return(commons.UserData{ID: 1, FullName: "Warga Konoha", PhoneNumber: "+628123456789"}, nil)

	mockRepository.EXPECT().CheckIfPhoneNumberExist(c.Request().Context(), "+628123456709").Return(false, nil)

	// fullname := "Kodok"
	// phoneNumber := "+628123456709"
	// payload := generated.UpdateProfileRequest{
	// 	FullName:    &fullname,
	// 	PhoneNumber: &phoneNumber,
	// }

	payload := repository.UpdateProfileInput{
		FullName:    "Kodok",
		PhoneNumber: "+628123456709",
	}

	mockRepository.EXPECT().UpdateProfile(c.Request().Context(), payload).Return(true, nil)

	server := NewServer(NewServerOptions{
		Repository: mockRepository,
		JWT:        mockJWT,
	})

	err := server.UpdateProfile(c)

	assert.Nil(t, err)
}

func TestServer_UpdateProfilePhoneNumberExist(t *testing.T) {
	token := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwiZnVsbF9uYW1lIjoiV2FyZ2EgS29ub2hhIiwicGhvbmVfbnVtYmVyIjoiNjI4MTIzNDU2Nzg5IiwiZXhwIjoxNjk0ODYxMzE4fQ.S7Ang3pIbMbU4NOPpb2360Qrxgs1kAE2Za20SOL_eC3hOaowpDfYtZ8JAaDFt-NdnU_S5D0yae-3s9n0GjwjmzBbcCzgQmNxRsP4HO1OZ0qJP688Wd1WQftuTapjblUrPPsgsaEiixUcDhLCDm5DE9vcUA-6BTYx3C7m3UeHxBhQX4oXqcFioMZg_qja3d2qSFlLjTEwnbuDL5mRaz2ac3cjB6Sv3xZf2aGGfrnPsnAKziX0wr3WN8H6QApi1SZLnYlrS89C5sIUyVSvVj4o5cU8u-2j6KXQfy3smYVBE4Hm5MGg4_395TjfIGdPf5fPc6IjPArCqpEkGlSMKIcXKIutW1Kv3WDhlAJ_WozN2aa6W4glNvcE3CSBqIqZbYzovncSEVEF-SePU8T8OI1zEjiwdmWX_jeCJMqUu-IdndHdtR9kf5N-m3SMui4N4F7hovjJNCh_0g1BjA0OaHJV62r1l7S1w0tY-vXkcglBLc_eRLFb6-zbhNGtKhKKDu2YgbfPoP5f8aivTWXNcXDV8NrXZuKXCuwhM1FU7lDj4FpYTerWzETXTZv3tzgSPAdO3coUwuhMlSVjMdb8bYzW5LHl6nQlATbj6FtNChCdaoymFp6nPm7li1IgqYluew1ahpoFCE8sFSbw2Z4__DHSSPiJR6JErFV3phuGk04CY-k"

	updateProfileJSON := `{"full_name": "Kodok", "phone_number": "+628123456709"}`

	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/users", strings.NewReader(updateProfileJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", token))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockCtrl := gomock.NewController(t)
	mockRepository := repository.NewMockRepositoryInterface(mockCtrl)

	mockJWT := commons.NewMockJWT(mockCtrl)

	mockJWT.EXPECT().ValidateToken(token).Return(commons.UserData{ID: 1, FullName: "Warga Konoha", PhoneNumber: "+628123456789"}, nil)

	mockRepository.EXPECT().CheckIfPhoneNumberExist(c.Request().Context(), "+628123456709").Return(true, nil)

	// fullname := "Kodok"
	// phoneNumber := "+628123456709"
	// payload := generated.UpdateProfileRequest{
	// 	FullName:    &fullname,
	// 	PhoneNumber: &phoneNumber,
	// }

	// payload := repository.UpdateProfileInput{
	// 	FullName:    "Kodok",
	// 	PhoneNumber: "+628123456709",
	// }

	// mockRepository.EXPECT().UpdateProfile(c.Request().Context(), payload).Return(true, nil)

	server := NewServer(NewServerOptions{
		Repository: mockRepository,
		JWT:        mockJWT,
	})

	err := server.UpdateProfile(c)

	assert.Nil(t, err)
}

func TestServer_UpdateProfilePhoneNumberNotValid(t *testing.T) {
	token := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwiZnVsbF9uYW1lIjoiV2FyZ2EgS29ub2hhIiwicGhvbmVfbnVtYmVyIjoiNjI4MTIzNDU2Nzg5IiwiZXhwIjoxNjk0ODYxMzE4fQ.S7Ang3pIbMbU4NOPpb2360Qrxgs1kAE2Za20SOL_eC3hOaowpDfYtZ8JAaDFt-NdnU_S5D0yae-3s9n0GjwjmzBbcCzgQmNxRsP4HO1OZ0qJP688Wd1WQftuTapjblUrPPsgsaEiixUcDhLCDm5DE9vcUA-6BTYx3C7m3UeHxBhQX4oXqcFioMZg_qja3d2qSFlLjTEwnbuDL5mRaz2ac3cjB6Sv3xZf2aGGfrnPsnAKziX0wr3WN8H6QApi1SZLnYlrS89C5sIUyVSvVj4o5cU8u-2j6KXQfy3smYVBE4Hm5MGg4_395TjfIGdPf5fPc6IjPArCqpEkGlSMKIcXKIutW1Kv3WDhlAJ_WozN2aa6W4glNvcE3CSBqIqZbYzovncSEVEF-SePU8T8OI1zEjiwdmWX_jeCJMqUu-IdndHdtR9kf5N-m3SMui4N4F7hovjJNCh_0g1BjA0OaHJV62r1l7S1w0tY-vXkcglBLc_eRLFb6-zbhNGtKhKKDu2YgbfPoP5f8aivTWXNcXDV8NrXZuKXCuwhM1FU7lDj4FpYTerWzETXTZv3tzgSPAdO3coUwuhMlSVjMdb8bYzW5LHl6nQlATbj6FtNChCdaoymFp6nPm7li1IgqYluew1ahpoFCE8sFSbw2Z4__DHSSPiJR6JErFV3phuGk04CY-k"

	updateProfileJSON := `{"full_name": "Kodok", "phone_number": "x8123456709xxxxx"}`

	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/users", strings.NewReader(updateProfileJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", token))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockCtrl := gomock.NewController(t)
	mockRepository := repository.NewMockRepositoryInterface(mockCtrl)

	mockJWT := commons.NewMockJWT(mockCtrl)

	mockJWT.EXPECT().ValidateToken(token).Return(commons.UserData{ID: 1, FullName: "Warga Konoha", PhoneNumber: "+628123456789"}, nil)

	// mockRepository.EXPECT().CheckIfPhoneNumberExist(c.Request().Context(), "+628123456709xxxx").Return(true, nil)

	// fullname := "Kodok"
	// phoneNumber := "+628123456709"
	// payload := generated.UpdateProfileRequest{
	// 	FullName:    &fullname,
	// 	PhoneNumber: &phoneNumber,
	// }

	// payload := repository.UpdateProfileInput{
	// 	FullName:    "Kodok",
	// 	PhoneNumber: "+628123456709",
	// }

	// mockRepository.EXPECT().UpdateProfile(c.Request().Context(), payload).Return(true, nil)

	server := NewServer(NewServerOptions{
		Repository: mockRepository,
		JWT:        mockJWT,
	})

	err := server.UpdateProfile(c)

	assert.Nil(t, err)
}

func TestServer_UpdateProfileFullNameNotValid(t *testing.T) {
	token := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwiZnVsbF9uYW1lIjoiV2FyZ2EgS29ub2hhIiwicGhvbmVfbnVtYmVyIjoiNjI4MTIzNDU2Nzg5IiwiZXhwIjoxNjk0ODYxMzE4fQ.S7Ang3pIbMbU4NOPpb2360Qrxgs1kAE2Za20SOL_eC3hOaowpDfYtZ8JAaDFt-NdnU_S5D0yae-3s9n0GjwjmzBbcCzgQmNxRsP4HO1OZ0qJP688Wd1WQftuTapjblUrPPsgsaEiixUcDhLCDm5DE9vcUA-6BTYx3C7m3UeHxBhQX4oXqcFioMZg_qja3d2qSFlLjTEwnbuDL5mRaz2ac3cjB6Sv3xZf2aGGfrnPsnAKziX0wr3WN8H6QApi1SZLnYlrS89C5sIUyVSvVj4o5cU8u-2j6KXQfy3smYVBE4Hm5MGg4_395TjfIGdPf5fPc6IjPArCqpEkGlSMKIcXKIutW1Kv3WDhlAJ_WozN2aa6W4glNvcE3CSBqIqZbYzovncSEVEF-SePU8T8OI1zEjiwdmWX_jeCJMqUu-IdndHdtR9kf5N-m3SMui4N4F7hovjJNCh_0g1BjA0OaHJV62r1l7S1w0tY-vXkcglBLc_eRLFb6-zbhNGtKhKKDu2YgbfPoP5f8aivTWXNcXDV8NrXZuKXCuwhM1FU7lDj4FpYTerWzETXTZv3tzgSPAdO3coUwuhMlSVjMdb8bYzW5LHl6nQlATbj6FtNChCdaoymFp6nPm7li1IgqYluew1ahpoFCE8sFSbw2Z4__DHSSPiJR6JErFV3phuGk04CY-k"

	updateProfileJSON := `{"full_name": "sd", "phone_number": "628123456789"}`

	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/users", strings.NewReader(updateProfileJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", token))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockCtrl := gomock.NewController(t)
	mockRepository := repository.NewMockRepositoryInterface(mockCtrl)

	mockJWT := commons.NewMockJWT(mockCtrl)

	mockJWT.EXPECT().ValidateToken(token).Return(commons.UserData{ID: 1, FullName: "Warga Konoha", PhoneNumber: "+628123456789"}, nil)

	// mockRepository.EXPECT().CheckIfPhoneNumberExist(c.Request().Context(), "+628123456709xxxx").Return(true, nil)

	// fullname := "Kodok"
	// phoneNumber := "+628123456709"
	// payload := generated.UpdateProfileRequest{
	// 	FullName:    &fullname,
	// 	PhoneNumber: &phoneNumber,
	// }

	// payload := repository.UpdateProfileInput{
	// 	FullName:    "Kodok",
	// 	PhoneNumber: "+628123456709",
	// }

	// mockRepository.EXPECT().UpdateProfile(c.Request().Context(), payload).Return(true, nil)

	server := NewServer(NewServerOptions{
		Repository: mockRepository,
		JWT:        mockJWT,
	})

	err := server.UpdateProfile(c)

	assert.Nil(t, err)
}
