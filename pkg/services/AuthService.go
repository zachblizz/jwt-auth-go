package services

import (
	"errors"
	"net/http"
	"time"
	
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"

	"bitbucket/zblizz/jwt-go/pkg/models"
	r "bitbucket/zblizz/jwt-go/pkg/repositories"
	"bitbucket/zblizz/jwt-go/pkg/utils"
)

// AuthService - the auth service struct
type AuthService struct {
	repository *r.UserRepository
}

// GetClaimsFromToken - checks if the user cookies are valid
func (s *AuthService) GetClaimsFromToken(tokenStr string) (*models.Claims, int, error) {
	claims := &models.Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return utils.GetJWTKey(), nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return claims, http.StatusUnauthorized, err
		}

		return claims, http.StatusBadRequest, err
	}

	if !token.Valid {
		return claims, http.StatusUnauthorized, errors.New("bad token")
	}

	return claims, http.StatusOK, nil
}

// IsAdminUser - checks the token if the user's role is admin or not
func (s *AuthService) IsAdminUser(req *http.Request, role string) int {
	// get the cookies
	c, err := req.Cookie(utils.TokenIdentifier)
	if err != nil {
		if err == http.ErrNoCookie {
			return http.StatusUnauthorized
		}

		return http.StatusBadRequest
	}

	claims, _, _ := s.GetClaimsFromToken(c.Value)

	if claims == nil || claims.Role != role {
		return http.StatusUnauthorized
	}

	return http.StatusOK
}

// RefreshTokens - create new auth and refresh tokens
func (s *AuthService) RefreshTokens(w *http.ResponseWriter, r *models.RefreshTokenReq) (int, error) {
	cookie := http.Cookie{
		Name:     utils.TokenIdentifier,
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
	}

	// extract the refresh token claim
	rClaims, rStatus, rErr := s.GetClaimsFromToken(r.RefreshToken)
	if rErr != nil {
		http.SetCookie(*w, &cookie)
		return rStatus, rErr
	}

	// ensure that the refresh token is expired
	elapsedTime := time.Unix(rClaims.ExpiresAt, 0).Sub(time.Now())
	if elapsedTime > utils.RefreshExpiration*time.Hour {
		http.SetCookie(*w, &cookie)
		return http.StatusForbidden, errors.New("refresh token expired")
	}

	expirationTime := time.Now().Add(utils.TokenExpiration * time.Minute)

	// get the user
	user, _ := s.repository.GetUserByID(rClaims.Subject)
	tokenStr, tErr := createToken(&user, expirationTime)
	if tErr != nil {
		http.SetCookie(*w, &cookie)
		return http.StatusInternalServerError, tErr
	}

	rt, rtErr := s.createRefreshToken(user.Username)
	if rtErr != nil {
		return http.StatusInternalServerError, rtErr
	}
	r.RefreshToken = rt

	cookie.Value = tokenStr
	cookie.Expires = expirationTime

	http.SetCookie(*w, &cookie)
	return http.StatusOK, nil
}

// CreateAuthToken - creates the authorization token and adds it to the cookie
func (s *AuthService) CreateAuthToken(userInfo *models.User, authResp *models.AuthorizeResp, cookie *http.Cookie) int {
	// create expirationTime of the token here
	expirationTime := time.Now().Add(utils.TokenExpiration * time.Minute)

	tokenStr, tErr := createToken(userInfo, expirationTime)
	if tErr != nil {
		authResp.Message = tErr.Error()
	}

	rt, rtErr := s.createRefreshToken(userInfo.Username)
	if rtErr != nil {
		authResp.Message = rtErr.Error()
	}

	// populate the cookie
	cookie.Value = tokenStr
	cookie.Expires = expirationTime

	authResp.Error = false
	authResp.IsAuthroized = true
	authResp.Role = userInfo.Role
	authResp.Message = ""
	authResp.RefreshToken = rt

	return http.StatusOK
}

func createToken(userInfo *models.User, expirationTime time.Time) (string, error) {
	// create the JWT claims, which includes the username, expiry time and user role
	claims := &models.Claims{
		Username: userInfo.Username,
		Role:     userInfo.Role,
		StandardClaims: jwt.StandardClaims{
			// in JWT, the expiry time is in milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// create the JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// create the JWT string
	tokenStr, err := token.SignedString(utils.GetJWTKey())
	if err != nil {
		return "", errors.New("could not create JWT token string")
	}

	return tokenStr, nil
}

func (s *AuthService) createRefreshToken(uname string) (string, error) {
	user, err := s.repository.GetUserByName(&models.Credentials{Username: uname})
	if err != nil {
		return "", err
	}

	expirationTime := time.Now().Add(utils.RefreshExpiration * time.Hour)

	claims := &models.Claims{
		StandardClaims: jwt.StandardClaims{
			Subject: user.ID.Hex(),
			// in JWT, the expiry time is in milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	rt, err := refreshToken.SignedString(utils.GetJWTKey())
	if err != nil {
		return "", err
	}

	return rt, nil
}

// HashPassword - hashes the password
func (s *AuthService) HashPassword(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	utils.Check(err)

	return string(hash)
}

// ComparePassword - compares the passwords
func (s *AuthService) ComparePassword(password, hash []byte) (bool, error) {
	err := bcrypt.CompareHashAndPassword(hash, password)
	return err == nil, err
}

// NewAuthService - creates a new auth service
func NewAuthService(repository *r.UserRepository) *AuthService {
	return &AuthService{repository: repository}
}
