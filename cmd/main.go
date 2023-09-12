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
)

func main() {
	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	var server generated.ServerInterface = newServer()

	// // Configure middleware with the custom claims type
	// config := echojwt.Config{
	// 	NewClaimsFunc: func(c echo.Context) jwt.Claims {
	// 		return new(commons.JWTCustomClaims)
	// 	},
	// 	SigningKey: []byte("secret"),
	// }
	// e.Use(echojwt.WithConfig(config))

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
