package services

import (
	"bitbucket/zblizz/jwt-go/models"
	r "bitbucket/zblizz/jwt-go/repositories"
	"bitbucket/zblizz/jwt-go/utils"
	"net/http"
	"time"
)

// UserService - service used to get access to user repo
type UserService struct {
	repository  *r.UserRepository
	authService *AuthService
}

// Signin - checks the db for the username and password
func (s *UserService) Signin(w *http.ResponseWriter, creds *models.Credentials) (models.AuthorizeResp, int) {
	cookie := http.Cookie{
		Name:     utils.TokenIdentifier,
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
	}

	authResp := models.AuthorizeResp{
		IsAuthroized: false,
		Error:        true,
		Role:         "",
		Message:      "credentials are invalid...",
		RefreshToken: "",
	}

	dbUser, err := s.repository.GetUserByName(creds)

	if err != nil {
		// remove the cookie...
		http.SetCookie(*w, &cookie)
		authResp.Message = err.Error()
		return authResp, http.StatusInternalServerError
	}

	goodPwd, _ := s.authService.ComparePassword([]byte(creds.Password), []byte(dbUser.Password))

	// check for bad credentials
	if !goodPwd {
		// remove the cookie...
		http.SetCookie(*w, &cookie)
		return authResp, http.StatusUnauthorized
	}

	userInfo := models.User{
		Username: creds.Username,
		Role:     dbUser.Role,
	}

	s.authService.CreateAuthToken(&userInfo, &authResp, &cookie)

	http.SetCookie(*w, &cookie)
	return authResp, http.StatusOK
}

// Signup - sign user up
func (s *UserService) Signup(reqUser *models.User) models.SignupResp {
	signupResp := models.SignupResp{
		Error:        false,
		ErrorMessage: "",
		SignupStatus: true,
	}

	reqUser.Password = s.authService.HashPassword([]byte(reqUser.Password))
	reqUser.Accepted = false

	err := s.repository.InsertUser(reqUser, "user-requests")
	if err != nil {
		signupResp.Error = true
		signupResp.ErrorMessage = err.Error()
		signupResp.SignupStatus = false
	}

	return signupResp
}

// AcceptUserRequest - create user from request
func (s *UserService) AcceptUserRequest(reqUser *models.User) models.SignupResp {
	signupResp := models.SignupResp{
		Error:        false,
		ErrorMessage: "",
		SignupStatus: true,
	}

	reqUser.Password = s.authService.HashPassword([]byte(reqUser.Password))

	err := s.repository.InsertUser(reqUser, "users")
	if err != nil {
		signupResp.Error = true
		signupResp.ErrorMessage = err.Error()
		signupResp.SignupStatus = false
	}

	signupResp, _ = s.repository.UpdateUserRequest(reqUser)

	return signupResp
}

// DeleteUser - deletes the user by the user id...
func (s *UserService) DeleteUser(user *models.User) error {
	return s.repository.DeleteUserByID(user.ID)
}

// GetUsers - gets the users by the collection
func (s *UserService) GetUsers(collection string) ([]models.User, error) {
	return s.repository.GetUsers(collection)
}

// NewUserService - creates a new user service
func NewUserService(repository *r.UserRepository, authService *AuthService) *UserService {
	return &UserService{repository: repository, authService: authService}
}
