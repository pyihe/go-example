package main

import (
	"context"

	"pkg"

	"github.com/pyihe/go-example/mongodb/service"
	"github.com/pyihe/go-example/mongodb/service/repository/mongo"
	"github.com/pyihe/go-pkg/tools"
	"github.com/pyihe/plogs"
	mgo "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const uri = ""

func main() {
	logger := pkg.InitLogger("mongo")
	defer logger.Close()

	mClient, err := mgo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		plogs.Fatalf("初始化MongoDB客户端发生错误: %v", err)
		return
	}
	defer mClient.Disconnect(context.Background())

	database := mClient.Database("example")

	userRepo := mongo.NewUserRepository(database)

	userService := service.NewUserService(userRepo)

	userService.Work()

	tools.Wait()
}
