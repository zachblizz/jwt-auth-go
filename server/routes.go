package server

import "net/http"

// Route - the route structure
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// GetRoutes - gets the list of routes for the site
func GetRoutes(s *Server) []Route {
	return []Route{
		Route{"test", "GET", "/test", s.test},
		Route{"signup", "PUT", "/signup", s.signup},
		Route{"signin", "POST", "/signin", s.signin},
		Route{"refresh", "POST", "/refresh", s.refreshAuthToken},
		Route{"logout", "POST", "/logout", s.logout},

		// admin Routes
		Route{"users", "GET", "/admin/users", s.getUsers},
		Route{"acceptUser", "POST", "/admin/acceptUser", s.acceptUser},
		Route{"deleteUser", "DELETE", "/admin/deleteUser", s.deleteUser},
	}
}
