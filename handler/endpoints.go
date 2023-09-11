package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/triardn/panganggo/generated"
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
	return nil
}

func (s *Server) GetUserDetailByID(ctx echo.Context, id int) error {
	return nil
}

func (s *Server) UpdateUser(ctx echo.Context) error {
	return nil
}
