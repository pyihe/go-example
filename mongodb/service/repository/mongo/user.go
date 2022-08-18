package mongo

import (
	"github.com/pyihe/go-example/mongodb/model"
	"github.com/pyihe/go-example/mongodb/service"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) service.UserAPI {
	res := &userRepository{}
	res.collection = db.Collection("user")

	// 创建索引

	return res
}

func (u *userRepository) InsertOne(user *model.User) error {
	return nil
}
