package main

import (
	"fmt"

	"go.uber.org/dig"
	mgo "gopkg.in/mgo.v2"

	r "bitbucket/zblizz/jwt-go/pkg/repositories"
	svr "bitbucket/zblizz/jwt-go/pkg/server"
	s "bitbucket/zblizz/jwt-go/pkg/services"
	"bitbucket/zblizz/jwt-go/pkg/utils"
)

func connectDatabase(config *utils.Config) (*mgo.Session, error) {
	fmt.Println("getting db connection")
	session, err := mgo.Dial(config.ConnString)

	if err == nil {
		session.SetMode(mgo.Monotonic, true)
	}

	fmt.Println("got db connection")
	return session, err
}

func buildContainer() *dig.Container {
	container := dig.New()

	container.Provide(utils.NewConfig)
	container.Provide(connectDatabase)
	container.Provide(r.NewUserRepository)
	container.Provide(s.NewUserService)
	container.Provide(s.NewAuthService)
	container.Provide(svr.NewServer)

	return container
}

func main() {
	container := buildContainer()
	utils.LoadProps()

	// TODO: might want to try to use wire for this
	// REF: https://github.com/google/wire
	err := container.Invoke(func(server *svr.Server) {
		server.Run()
	})

	if err != nil {
		panic(err)
	}
}
