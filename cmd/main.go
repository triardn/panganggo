package main

import (
	"log"
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

	privateKey, err := os.ReadFile("cert/id_rsa")
	if err != nil {
		log.Fatalln(err)
	}
	publicKey, err := os.ReadFile("cert/id_rsa.pub")
	if err != nil {
		log.Fatalln(err)
	}

	newJWT := commons.NewJWT(privateKey, publicKey)

	opts := handler.NewServerOptions{
		Repository: repo,
		JWT:        newJWT,
	}

	return handler.NewServer(opts)
}
