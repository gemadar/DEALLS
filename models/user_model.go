package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User model
type User struct {
	Id       primitive.ObjectID `json:"id,omitempty"`
	Username string             `json:"username,omitempty" validate:"required"`
	Password string             `json:"password,omitempty" validate:"required"`
	Name     string             `json:"name,omitempty" validate:"required"`
	Location string             `json:"location,omitempty" validate:"required"`
	Age      int                `json:"age,omitempty" validate:"required"`
	Gender   string             `json:"gender,omitempty" validate:"required"`
	Viewed   int                `json:"viewed,omitempty"`
	Premium  bool               `json:"premium,omitempty" default:"false"`
	Liked    []string           `json:"liked,omitempty"`
	Passed   []string           `json:"passed,omitempty"`
}
