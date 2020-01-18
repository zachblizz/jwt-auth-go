package repositories

import (
	"bitbucket/zblizz/jwt-go/pkg/models"
	"errors"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// UserRepository - repo to gather user info from mongo
type UserRepository struct {
	session *mgo.Session
}

const dbName = "auth"

// GetUsers - gets the users from the collection provided
func (repo *UserRepository) GetUsers(collection string) ([]models.User, error) {
	var err error
	c := repo.session.DB(dbName).C(collection)

	var results []models.User
	if collection == "user-requests" {
		err = c.Find(bson.M{"accepted": false}).All(&results)
	} else {
		err = c.Find(nil).All(&results)
	}

	return results, err
}

// GetUserByName - retrieves the user from the DB based on the username provided
func (repo *UserRepository) GetUserByName(creds *models.Credentials) (models.User, error) {
	c := repo.session.DB(dbName).C("users")
	var dbUser models.User
	err := c.Find(bson.M{"username": creds.Username}).One(&dbUser)

	return dbUser, err
}

// GetUserByID - retrieves the user by ID
func (repo *UserRepository) GetUserByID(id string) (models.User, error) {
	c := repo.session.DB(dbName).C("users")
	var dbUser models.User
	err := c.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&dbUser)

	return dbUser, err
}

// InsertUser - inserts the user if not already in DB
func (repo *UserRepository) InsertUser(reqUser *models.User, collection string) error {
	c := repo.session.DB(dbName).C(collection)
	count, countErr := c.Find(bson.M{"username": reqUser.Username}).Limit(1).Count()

	if countErr != nil {
		return countErr
	}

	if count > 0 {
		return errors.New("user pending acceptance")
	}

	err := c.Insert(reqUser)
	return err
}

// AcceptUserRequest - create user from request
func (repo *UserRepository) AcceptUserRequest(reqUser *models.User) models.SignupResp {
	signupResp := models.SignupResp{
		ErrorMessage: "",
		Error:        false,
		SignupStatus: true,
	}

	err := repo.InsertUser(reqUser, "users")
	if err != nil {
		signupResp.Error = true
		signupResp.ErrorMessage = err.Error()
		signupResp.SignupStatus = false
	}

	signupResp, _ = repo.UpdateUserRequest(reqUser)

	return signupResp
}

// UpdateUserRequest - updates the user request accepted variable to true
func (repo *UserRepository) UpdateUserRequest(user *models.User) (models.SignupResp, error) {
	resp := models.SignupResp{
		Error:        false,
		ErrorMessage: "",
		SignupStatus: true,
	}

	c := repo.session.DB(dbName).C("user-requests")
	err := c.UpdateId(user.ID, bson.M{"$set": bson.M{"accepted": true}})
	if err != nil {
		resp.Error = true
		resp.ErrorMessage = err.Error()
		resp.SignupStatus = false
	}

	return resp, err
}

// DeleteUserByID - deletes a user by ID
func (repo *UserRepository) DeleteUserByID(id bson.ObjectId) error {
	c := repo.session.DB(dbName).C("users")
	return c.RemoveId(id)
}

// NewUserRepository - creates a new user repository
func NewUserRepository(session *mgo.Session) *UserRepository {
	return &UserRepository{session: session}
}
