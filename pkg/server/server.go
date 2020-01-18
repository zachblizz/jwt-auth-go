package server

import (
	s "bitbucket/zblizz/jwt-go/pkg/services"
	utils "bitbucket/zblizz/jwt-go/pkg/utils"
	"log"
	"net/http"

	"github.com/rs/cors"
)

// Server - server struct
type Server struct {
	config      *utils.Config
	userService *s.UserService
	authService *s.AuthService
}

// Run - runs the server
func (s *Server) Run() {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000", "http://localhost:4000"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	})

	httpServer := &http.Server{
		Addr:    s.config.Port,
		Handler: c.Handler(NewRouter(s)),
	}

	log.Fatal(httpServer.ListenAndServe())
}

// NewServer - creates a new server with DI
func NewServer(userService *s.UserService, authService *s.AuthService, config *utils.Config) *Server {
	return &Server{
		config:      config,
		userService: userService,
		authService: authService,
	}
}
