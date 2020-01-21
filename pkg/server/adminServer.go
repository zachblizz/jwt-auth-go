package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"bitbucket/zblizz/jwt-go/pkg/models"
	utils "bitbucket/zblizz/jwt-go/pkg/utils"
)

func adminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func (s *Server) getUsers(w http.ResponseWriter, req *http.Request) {
	collection, ok := req.URL.Query()["collection"]
	isAdminStatus := s.authService.IsAdminUser(req, "admin")

	if ok && len(collection) == 1 && isAdminStatus == http.StatusOK {
		users, err := s.userService.GetUsers(collection[0])
		utils.CheckAndWriteHeader(err, w, http.StatusInternalServerError)
		usersResp := models.UsersResp{
			Count: len(users),
			Users: users,
		}

		resp, err := json.MarshalIndent(usersResp, "", "  ")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, string(resp))
		return
	}

	w.WriteHeader(isAdminStatus)
	fmt.Fprintf(w, "Nope!")
}

func (s *Server) acceptUser(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var userReq models.User
	err := decoder.Decode(&userReq)
	utils.CheckAndWriteHeader(err, w, http.StatusInternalServerError)

	isAdminStatus := s.authService.IsAdminUser(req, "admin")

	if isAdminStatus == http.StatusOK {
		rawResp := s.userService.AcceptUserRequest(&userReq)
		resp, _ := json.MarshalIndent(rawResp, "", "  ")
		fmt.Fprintf(w, string(resp))
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Nope!")
	}
}

func (s *Server) deleteUser(w http.ResponseWriter, req *http.Request) {
	isAdminStatus := s.authService.IsAdminUser(req, "admin")
	msg := "Deleted!"

	if isAdminStatus == http.StatusOK {
		decoder := json.NewDecoder(req.Body)
		var userInfo models.User
		err := decoder.Decode(&userInfo)
		utils.CheckAndWriteHeader(err, w, http.StatusInternalServerError)

		err = s.userService.DeleteUser(&userInfo)
		utils.CheckAndWriteHeader(err, w, http.StatusInternalServerError)
	} else {
		msg = "Nope!"
	}

	w.WriteHeader(isAdminStatus)
	fmt.Fprintf(w, msg)
}
