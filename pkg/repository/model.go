package repository

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID          primitive.ObjectID `bson:"ID"`
	FistName    string             `bson:"FistName"`
	LastName    string             `bson:"LastName"`
	Birthday    time.Time          `bson:"Birthday"`
	Email       string             `bson:"Email"`
	Username    string             `bson:"Username"`
	Password    string             `bson:"Password"`
	PhoneNumber string             `bson:"PhoneNumber"`
	Photo       []byte             `bson:"Photo"`
	AvailableMb uint32             `bson:"AvailableMb"`
	CreatedAt   time.Time          `bson:"CreatedAt"`
	LastLogin   time.Time          `bson:"LastLogin"`
	LastLogout  time.Time          `bson:"LastLogout"`
	LastEdited  time.Time          `bson:"LastEdited"`
	DeletedAt   time.Time          `bson:"DeletedAt"`
	Deleted     bool               `bson:"Deleted"`
	Activated   bool               `bson:"Activated"`
	Locked      bool               `bson:"Locked"`
	Salt        []byte             `bson:"Salt"`
}

type Payement struct {
	ID         primitive.ObjectID `bson:"ID"`
	PlanID     primitive.ObjectID `bson:"PlanID"`
	UserID     primitive.ObjectID `bson:"UserID"`
	ExecutedAt time.Time          `bson:"ExecutedAt"`
}

type Plan struct {
	ID           primitive.ObjectID `bson:"ID"`
	DurationDays uint32             `bson:"DurationDays"`
	AvailableMB  uint32             `bson:"AvailableMB"`
	PriceRON     uint32             `bson:"PriceRon"`
}
