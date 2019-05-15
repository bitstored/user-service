package service

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Activation struct {
	ID        primitive.ObjectID `bson:"ID"`
	Token     []byte             `bson:"Token"`
	ExpiresAt time.Time          `bson:"ExpiresAt"`
}

type Session struct {
	ID        primitive.ObjectID `bson:"ID"`
	Token     []byte
	ExpiresAt time.Time
}
