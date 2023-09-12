package handler

import (
	"github.com/triardn/panganggo/commons"
	"github.com/triardn/panganggo/repository"
)

type Server struct {
	Repository repository.RepositoryInterface
	JWT        commons.JWT
}

type NewServerOptions struct {
	Repository repository.RepositoryInterface
	JWT        commons.JWT
}

func NewServer(opts NewServerOptions) *Server {
	return &Server{
		Repository: opts.Repository,
		JWT:        opts.JWT,
	}
}
