package utils

import "fmt"

// Config - the config for the server
type Config struct {
	Port       string
	ConnString string
}

// NewConfig - creates a new config object
func NewConfig() *Config {
	props := GetProps()
	host, _ := props.GetString("db", "host")
	port, _ := props.GetString("db", "port")

	return &Config{
		Port:       ":4000",
		ConnString: fmt.Sprintf("%s:%s", host, port),
	}
}
