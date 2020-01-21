package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"bitbucket/zblizz/jwt-go/pkg/models"
	utils "bitbucket/zblizz/jwt-go/pkg/utils"
)

func (s *Server) test(w http.ResponseWriter, req *http.Request) {
	resp, err := json.MarshalIndent("Hello, World!", "", "  ")
	utils.CheckAndWriteHeader(err, w, http.StatusInternalServerError)
	if err != nil {
		return
	}

	fmt.Fprintf(w, string(resp))
}

func (s *Server) signin(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var authReq models.Credentials
	err := decoder.Decode(&authReq)
	utils.CheckAndWriteHeader(err, w, http.StatusInternalServerError)
	if err != nil {
		return
	}

	authResp, status := s.userService.Signin(&w, &authReq)
	resp, err := json.MarshalIndent(authResp, "", "  ")
	utils.CheckAndWriteHeader(err, w, http.StatusInternalServerError)
	if err != nil {
		return
	}

	w.WriteHeader(status)
	fmt.Fprintf(w, string(resp))
}

func (s *Server) signup(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var signupReq models.User
	err := decoder.Decode(&signupReq)
	utils.CheckAndWriteHeader(err, w, http.StatusInternalServerError)
	if err != nil {
		return
	}

	signupResp := s.userService.Signup(&signupReq)
	resp, err := json.MarshalIndent(signupResp, "", "  ")
	utils.CheckAndWriteHeader(err, w, http.StatusInternalServerError)
	if err != nil {
		return
	}

	fmt.Fprintf(w, string(resp))
	w.WriteHeader(http.StatusOK)
}

func (s *Server) refreshAuthToken(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var refreshReq models.RefreshTokenReq
	err := decoder.Decode(&refreshReq)
	utils.CheckAndWriteHeader(err, w, http.StatusInternalServerError)
	if err != nil {
		return
	}

	status, rErr := s.authService.RefreshTokens(&w, &refreshReq)
	w.WriteHeader(status)
	if rErr != nil {
		fmt.Fprintf(w, rErr.Error())
		return
	}

	resp, err := json.MarshalIndent(refreshReq, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, string(resp))
}

func (s *Server) logout(w http.ResponseWriter, req *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     utils.TokenIdentifier,
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
	})

	resp, err := json.MarshalIndent(models.User{}, "", "  ")
	utils.CheckAndWriteHeader(err, w, http.StatusInternalServerError)
	if err != nil {
		return
	}

	// TODO: need to expire the user's refresh token

	fmt.Fprintf(w, string(resp))
}
