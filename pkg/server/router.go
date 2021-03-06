package server

import (
	"net/http"
	
	"github.com/gorilla/mux"

	"bitbucket/zblizz/jwt-go/pkg/utils"
)

// NewRouter - creates the router for the service
func NewRouter(s *Server) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range GetRoutes(s) {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = utils.Logger(handler, route.Name)

		router.
			PathPrefix("/api/v1").
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)

	}

	return router
}
