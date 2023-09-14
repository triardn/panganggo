package handler

import (
	"database/sql"
	"encoding/json"
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

var (
	registrationJSON = `{"full_name": "Petrikor", "phone_number": "+6281234567", "password": "petr1k0r@10ss"}`
	loginJSON        = `{"phone_number": "+6281234567", "password": "abcdefg"}`
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

	mockRepository.EXPECT().CheckIfPhoneNumberExist(c, "+6281234567").Return(false, sql.ErrConnDone)

	server := NewServer(NewServerOptions{
		Repository: mockRepository,
		JWT:        mockJWT,
	})

	generated.RegisterHandlers(e, server)

	server.Registration(c)
	// t.Errorf("Unexpected error: %v", err)
}

func TestLogin(t *testing.T) { // BROKEN UNIT TEST
	// // Setup
	// e := echo.New()
	// req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(loginJSON))
	// req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	// rec := httptest.NewRecorder()
	// c := e.NewContext(req, rec)

	// mockCtrl := gomock.NewController(t)
	// defer mockCtrl.Finish()
	// mockRepository := repository.NewMockRepositoryInterface(mockCtrl)

	// mockJWT := commons.NewMockJWT(mockCtrl)

	// mockRepository.EXPECT().GetUsersByPhoneNumber(c, "+6281234567").Return(repository.GetUsersByPhoneNumberOutput{}, sql.ErrConnDone)

	// server := NewServer(NewServerOptions{
	// 	Repository: mockRepository,
	// 	JWT:        mockJWT,
	// })

	// generated.RegisterHandlers(e, server)

	// _ = server.Registration(c)

	// req, err := http.NewRequest("POST", "/login", strings.NewReader(loginJSON))
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// // We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	// rr := httptest.NewRecorder()
	// handler := http.HandlerFunc(HealthCheckHandler)

	// server := NewServer(NewServerOptions{
	// 	Repository: mockRepository,
	// 	JWT:        mockJWT,
	// })

	// generated.RegisterHandlers(e, server)

	// // Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// // directly and pass in our Request and ResponseRecorder.
	// handler.ServeHTTP(rr, req)
	// // Check the status code is what we expect.
	// if status := rr.Code; status != http.StatusOK {
	// 	t.Errorf("handler returned wrong status code: got %v want %v",
	// 		status, http.StatusOK)
	// }
	// // Check the response body is what we expect.
	// expected := `{"alive": true}`
	// if rr.Body.String() != expected {
	// 	t.Errorf("handler returned unexpected body: got %v want %v",
	// 		rr.Body.String(), expected)
	// }

}

func TestServer_Login(t *testing.T) {
	// type fields struct {
	// 	Repository repository.RepositoryInterface
	// 	JWT        commons.JWT
	// }
	// type args struct {
	// 	ctx echo.Context
	// }
	// tests := []struct {
	// 	name    string
	// 	fields  fields
	// 	args    args
	// 	wantErr bool
	// }{
	// 	// TODO: Add test cases.
	// }
	// for _, tt := range tests {
	// 	t.Run(tt.name, func(t *testing.T) {
	// 		s := &Server{
	// 			Repository: tt.fields.Repository,
	// 			JWT:        tt.fields.JWT,
	// 		}
	// 		if err := s.Login(tt.args.ctx); (err != nil) != tt.wantErr {
	// 			t.Errorf("Server.Login() error = %v, wantErr %v", err, tt.wantErr)
	// 		}
	// 	})
	// }

	loginJSON := generated.LoginRequest{
		PhoneNumber: "+6281234567",
		Password:    "abcdefg",
	}

	b, _ := json.Marshal(loginJSON)

	e := echo.New()

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(string(b)))
	rec := httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	c := e.NewContext(req, rec)

	// u := new(generated.LoginRequest)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRepository := repository.NewMockRepositoryInterface(mockCtrl)

	mockJWT := commons.NewMockJWT(mockCtrl)

	mockRepository.EXPECT().GetUsersByPhoneNumber(c, "+6281234567").Return(repository.GetUsersByPhoneNumberOutput{}, sql.ErrConnDone)

	s := &Server{
		Repository: mockRepository,
		JWT:        mockJWT,
	}

	err := s.Login(c)
	if err != nil {
		t.Errorf("Server.Login() error = %v, wantErr %v", err, nil)
	}

	// err := c.Bind(&u)
	// fmt.Println(u)
	// if assert.NoError(t, err) {
	// 	assert.NotNil(t, u.PhoneNumber)
	// 	assert.NotNil(t, u.Password)
	// } else {
	// 	t.Errorf("error: %v", err)
	// }
}
