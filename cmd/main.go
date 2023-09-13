package main

import (
	"os"

	"github.com/triardn/panganggo/commons"
	"github.com/triardn/panganggo/generated"
	"github.com/triardn/panganggo/handler"
	"github.com/triardn/panganggo/repository"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/rs/zerolog"
)

func main() {
	e := echo.New()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := zerolog.New(os.Stdout)
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info().
				Str("URI", v.URI).
				Int("status", v.Status).
				Msg("request")

			return nil
		},
	}))

	var server generated.ServerInterface = newServer()

	generated.RegisterHandlers(e, server)
	e.Logger.Fatal(e.Start(":1323"))
}

func newServer() *handler.Server {
	dbDsn := os.Getenv("DATABASE_URL")
	var repo repository.RepositoryInterface = repository.NewRepository(repository.NewRepositoryOptions{
		Dsn: dbDsn,
	})

	privateKey := []byte(os.Getenv("ID_RSA"))
	publicKey := []byte(os.Getenv("ID_RSA_PUB"))

	newJWT := commons.NewJWT(privateKey, publicKey)

	opts := handler.NewServerOptions{
		Repository: repo,
		JWT:        newJWT,
	}

	return handler.NewServer(opts)
}
