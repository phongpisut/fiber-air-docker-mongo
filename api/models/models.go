package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id      primitive.ObjectID `json:"id,omitempty"`
	Name    string             `json:"name,omitempty" validate:"required"`
	Surname string             `json:"surname,omitempty"`
	Tel     string             `json:"tel,omitempty" validate:"required"`
}

type Properties struct {
	Id      primitive.ObjectID `json:"id,omitempty"`
	Maximum int                `json:"maximum,omitempty"`
}

type UserAndProperties struct {
	User     []User     `json:"user"`
	Settings Properties `json:"settings"`
}

type Error struct {
	Message string
}
