package models

import (
	"github.com/dgrijalva/jwt-go"
	"gopkg.in/mgo.v2/bson"
)

// User - request payload for the authorization
// and model for DB
type User struct {
	ID       bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Accepted bool          `bson:"accepted"`
	Role     string        `bson:"role" json:"role"`
	Password string        `bson:"password" json:"password"`
	Username string        `bson:"username" json:"username"`
	Email    string        `bson:"email" json:"email"`
}

// Credentials - request authorization payload
type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// AuthorizeResp - response payload
type AuthorizeResp struct {
	IsAuthroized bool   `json:"isAuthorized"`
	Error        bool   `json:"error"`
	Role         string `json:"role"`
	Message      string `json:"message"`
	RefreshToken string `json:"refreshToken"`
}

// SignupResp - signup resp
type SignupResp struct {
	SignupStatus bool   `json:"signupStatus"`
	Error        bool   `json:"error"`
	ErrorMessage string `json:"message"`
}

// Claims - model used to store the JWT information
type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

// UsersResp - model used to hold users information
type UsersResp struct {
	Count int    `json:"count"`
	Users []User `json:"users"`
}

// RefreshTokenReq - model used for refresh token post request
type RefreshTokenReq struct {
	RefreshToken string `json:"refreshToken"`
}
