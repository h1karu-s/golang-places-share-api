package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Users .
type Users struct {
	Users []User `json:"users"` 
}

//User .
type User struct{
	ID    primitive.ObjectID    `json:"id" bson:"_id,omitempty"`
	Name   string   `json:"name"`
	Email  string   `json:"email"`
	Image  string   `json:"image"`
	Places []string `json:"places"`
}

//UserContainPassword has +password
type UserContainPassword struct{
	ID       primitive.ObjectID    `json:"id" bson:"_id,omitempty"`
	Name     string                `json:"name"`
	Email    string                `json:"email"`
	Password string                `json:"password"`
	Image    string                `json:"image"`
	Places   []string              `json:"places"`
}

//SignUpUser .
type SignUpUser struct {
	UserID  primitive.ObjectID `json:"userId"`
	Email   string             `josn:"email"`
	Token   string             `json:"token"`
}