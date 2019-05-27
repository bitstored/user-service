package service

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Activation struct {
	Email     string    `bson:"Email"`
	Token     string    `bson:"Token"`
	ExpiresAt time.Time `bson:"ExpiresAt"`
}

type Session struct {
	ID primitive.ObjectID `bson:"ID"`

	FirstName string
	LastName  string
	Token     string
	ExpiresAt time.Time
}
