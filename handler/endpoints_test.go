package handler

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/triardn/panganggo/commons"
	"github.com/triardn/panganggo/generated"
	"github.com/triardn/panganggo/repository"
	"go.uber.org/mock/gomock"
)

func TestHello(t *testing.T) {

}

var (
	registrationJSON = `{"full_name": "Petrikor", "phone_number": "+6281234567", "password": "petr1k0r@10ss"}`
)

func TestRegister(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/registration", strings.NewReader(registrationJSON))
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

	generated.RegisterHandlers(e, server)

	mockRepository.EXPECT().CheckIfPhoneNumberExist(c, "+6281234567").Return(false, sql.ErrConnDone)

	server.Registration(c)
}
