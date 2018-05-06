package data

import "github.com/globalsign/mgo/bson"

type User struct {
	ID       bson.ObjectId `json:"id" bson:"_id"`
	User     string        `json:"user"`
	Password string        `json:"password"`
}

type Users []User
