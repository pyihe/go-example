package model

import "time"

type User struct {
	ID         string              `bson:"id"            json:"id"`
	Name       string              `bson:"name"          json:"name"`
	Phone      string              `bson:"phone"         json:"phone"`
	Age        int                 `bson:"age"           json:"age"`
	Sex        uint8               `bson:"sex"           json:"sex"`
	Department int                 `bson:"department"    json:"department"`
	Favorite   []*Favorite         `bson:"favorite"      json:"favorite"`
	Idol       map[string]struct{} `bson:"idol"          json:"idol"`
	CreatedAt  *time.Time          `bson:"created_at"    json:"created_at"`
}

type Favorite struct {
	Name  string `bson:"name"          json:"name"`
	Start string `bson:"start"         json:"start"`
}
