package utils

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/dlintw/goconf"
)

// TokenIdentifier - the id for the token used in the cookie
const TokenIdentifier = "token"

// TokenExpiration - the time in minutes the token will expire
const TokenExpiration = 15

// RefreshExpiration - the time in hours the refresh token will expire
const RefreshExpiration = 24

// CheckAndWriteHeader - check the error being passed in
func CheckAndWriteHeader(e error, w http.ResponseWriter, status int) {
	if err != nil {
		w.WriteHeader(status)
		fmt.Fprintf(w, err.Error())
	}
}

// Check - panic the error
func Check(e error) {
	if e != nil {
		panic(e)
	}
}

var props *goconf.ConfigFile
var err error

// LoadProps - loads the props for the app
func LoadProps() {
	if props == nil {
		fmt.Println("loading properties")

		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		fmt.Println(dir)
		if err != nil {
			panic(err)
		}
		props, err = goconf.ReadConfigFile(filepath.Join(dir, "props.cfg"))
		if err != nil {
			panic(err)
		}

		fmt.Println("done loading props")
	}
}

// GetProps - gets the properties for the application
func GetProps() *goconf.ConfigFile {
	if err == nil {
		return props
	}

	return nil
}

// GetJWTKey - returns the byte array containing the JWT key
func GetJWTKey() []byte {
	key, _ := props.GetString("jwt", "key")
	return []byte(key)
}
