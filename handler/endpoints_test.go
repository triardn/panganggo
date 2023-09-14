package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/triardn/panganggo/commons"
	"github.com/triardn/panganggo/repository"
	"go.uber.org/mock/gomock"
)

func TestHello(t *testing.T) {

}

var (
	registrationJSON = `{"full_name": "Petrikor", "phone_number": "+6281234567", "password": "Petr1k0r@10ss"}`
)

func TestRegister(t *testing.T) { // BROKEN UNIT TEST
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/registration", strings.NewReader(registrationJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockCtrl := gomock.NewController(t)
	mockRepository := repository.NewMockRepositoryInterface(mockCtrl)

	mockJWT := commons.NewMockJWT(mockCtrl)

	mockRepository.EXPECT().CheckIfPhoneNumberExist(c.Request().Context(), "+6281234567").Return(false, nil)

	// create new user
	payload := repository.RegisterInput{
		FullName:    "Petrikor",
		PhoneNumber: "+6281234567",
		Password:    commons.HashAndSalt([]byte("Petr1k0r@10ss")),
	}

	mockRepository.EXPECT().Register(c.Request().Context(), payload).Return(repository.RegisterOutput{}, nil)

	server := NewServer(NewServerOptions{
		Repository: mockRepository,
		JWT:        mockJWT,
	})

	_ = server.Registration(c)
}
