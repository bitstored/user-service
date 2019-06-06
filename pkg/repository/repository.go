package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	ID          = "ID"
	FIRSTNAME   = "FirstName"
	LASTNAME    = "LastName"
	BIRTHDAY    = "Birthday"
	EMAIL       = "Email"
	USERNAME    = "Username"
	PASSWORD    = "Password"
	PHONENUMBER = "PhoneNumber"
	PHOTO       = "Photo"
	AVAILABLEMB = "AvailableMb"
	CRAETEDAT   = "CreatedAt"
	LASTLOGIN   = "LastLogin"
	LASTLOGOUT  = "LastLogout"
	LASTEDITED  = "LastEdited"
	DELETEDAT   = "DeletedAt"
	ACTIVATED   = "Activated"
	LOCKED      = "Locked"
	DELETED     = "Deleted"
	SALT        = "Salt"
	DefaultMb   = 50 * 1024
)

type Repository struct {
	*mongo.Client
	DatabaseName string
	Lock         sync.RWMutex
}

// "mongodb://localhost:27017"
func NewRepository(host string, dbName string) (*Repository, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(host))
	if err != nil {
		return nil, err
	}

	// Check the connection
	err = client.Ping(context.Background(), readpref.Primary())

	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to MongoDB!")
	return &Repository{Client: client, DatabaseName: dbName}, nil
}

func (r *Repository) CreateAccount(ctx context.Context, collectionName string, user User) (*mongo.InsertOneResult, error) {
	database := r.Database(r.DatabaseName)
	collection := database.Collection(collectionName)

	r.Lock.RLock()
	if !r.dataExists(ctx, collection, EMAIL, user.Email) {
		err := fmt.Errorf("account can't be created, %s = %s is already in use", EMAIL, user.Email)
		LogEvent(EVENT_TYPE_INSERT, STATUS_ERROR, err, user)
		return nil, err
	}
	if !r.dataExists(ctx, collection, USERNAME, user.Username) {
		err := fmt.Errorf("account can't be created, %s = %s is already in use", USERNAME, user.Username)
		LogEvent(EVENT_TYPE_INSERT, STATUS_ERROR, err, user)
		return nil, err
	}
	if !r.dataExists(ctx, collection, PHONENUMBER, user.PhoneNumber) {
		err := fmt.Errorf("account can't be created, %s = %s is already in use", PHONENUMBER, user.PhoneNumber)
		LogEvent(EVENT_TYPE_INSERT, STATUS_ERROR, err, user)
		return nil, err
	}
	r.Lock.RUnlock()

	user.AvailableMb = DefaultMb
	user.Activated = true // TODO
	user.Locked = false
	user.Deleted = false
	user.CreatedAt = time.Now()

	r.Lock.Lock()
	res, err := collection.InsertOne(ctx, user)
	r.Lock.Unlock()

	if err != nil {
		LogEvent(EVENT_TYPE_INSERT, STATUS_ERROR, err, user)
		return nil, err
	}
	if res.InsertedID != nil {
		LogEvent(EVENT_TYPE_INSERT, STATUS_SUCCESS, nil, user)
	}
	return res, nil
}

func (r *Repository) ActivateAccount(ctx context.Context, collectionName string, email string) (*mongo.UpdateResult, error) {
	database := r.Database(r.DatabaseName)
	collection := database.Collection(collectionName)
	filter := bson.D{
		{EMAIL, email},
		// {DELETED, false},
		// {ACTIVATED, false},
		// {LOCKED, false},
	}
	update := bson.D{
		{"$set", bson.D{
			{ACTIVATED, true},
			{LASTEDITED, time.Now().String()},
		}},
	}
	res, err := collection.UpdateMany(ctx, filter, update)
	if err != nil {
		LogEvent(EVENT_TYPE_UPDATE, STATUS_ERROR, err, filter)
		return nil, err
	}
	if res.ModifiedCount > 0 {
		LogEvent(EVENT_TYPE_UPDATE, STATUS_SUCCESS, nil, filter)

	}
	LogEvent(EVENT_TYPE_UPDATE, STATUS_ERROR, err, filter)
	return res, nil
}

func (r *Repository) UpdateAccount(ctx context.Context, collectionName string, userID primitive.ObjectID, user User) (*mongo.UpdateResult, error) {
	database := r.Database(r.DatabaseName)
	collection := database.Collection(collectionName)
	filter := bson.D{
		{ID, userID},
		{EMAIL, user.Email},
		{USERNAME, user.Username},
		{PASSWORD, user.Password},
		{SALT, user.Salt},
		{DELETED, false},
		{ACTIVATED, true},
		{LOCKED, false},
	}
	update := bson.D{
		{"$set", bson.D{
			{FIRSTNAME, user.FirstName},
			{LASTNAME, user.LastName},
			{PHONENUMBER, user.PhoneNumber},
			{PHOTO, user.Photo},
			{BIRTHDAY, user.Birthday},
			{LASTEDITED, time.Now()},
		}},
	}
	res, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	if res.ModifiedCount == 0 {
		return nil, fmt.Errorf("user not found, account may be locked, deleted or not activated")
	}
	return res, nil
}

func (r *Repository) DeleteAccount(ctx context.Context, collectionName string, userID primitive.ObjectID, password string) (*mongo.UpdateResult, error) {
	database := r.Database(r.DatabaseName)
	collection := database.Collection(collectionName)
	filter := bson.D{
		{ID, userID},
		{DELETED, false},
		{ACTIVATED, true},
		{LOCKED, false},
		{PASSWORD, password},
	}
	update := bson.D{
		{"$set", bson.D{
			{DELETEDAT, time.Now()},
			{DELETED, true},
		}},
	}
	res, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	if res.ModifiedCount == 0 {
		return nil, fmt.Errorf("account is either inactive or inexistent")
	}
	return res, nil
}

func (r *Repository) GetAccount(ctx context.Context, collectionName string, userID primitive.ObjectID) *User {
	database := r.Database(r.DatabaseName)
	collection := database.Collection(collectionName)
	filter := bson.D{
		{ID, userID},
	}

	res := collection.FindOne(ctx, filter)

	user := new(User)
	res.Decode(&user)

	return user
}

func (r *Repository) Login(ctx context.Context, collectionName, username, password string) *User {
	database := r.Database(r.DatabaseName)
	collection := database.Collection(collectionName)
	filter := bson.D{
		{USERNAME, username},
		{DELETED, false},
		{ACTIVATED, true},
		{LOCKED, false},
	}
	update := bson.D{
		{"$set", bson.D{
			{LASTLOGIN, time.Now()},
		}},
	}
	res := collection.FindOne(ctx, filter)
	ures, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil
	}
	if ures.ModifiedCount == 0 {
		return nil
	}
	user := new(User)
	res.Decode(&user)

	return user
}

func (r *Repository) ResetPassword(ctx context.Context, collectionName, userID, password string, newPassword string) (*mongo.UpdateResult, error) {
	database := r.Database(r.DatabaseName)
	collection := database.Collection(collectionName)
	filter := bson.D{
		{ID, userID},
		{PASSWORD, password},
		{DELETED, false},
		{ACTIVATED, true},
		{LOCKED, false},
	}
	update := bson.D{
		{"$set", bson.D{
			{LASTEDITED, time.Now()},
			{PASSWORD, newPassword},
		}},
	}
	res, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *Repository) LockAccount(ctx context.Context, collectionName, userID, password string) (*mongo.UpdateResult, error) {
	database := r.Database(r.DatabaseName)
	collection := database.Collection(collectionName)
	filter := bson.D{
		{ID, userID},
		{PASSWORD, password},
		{DELETED, false},
		{ACTIVATED, true},
		{LOCKED, false},
	}
	update := bson.D{
		{"$set", bson.D{
			{LASTEDITED, time.Now()},
			{LOCKED, true},
		}},
	}
	res, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *Repository) UnlockAccount(ctx context.Context, collectionName, userID, password string) (*mongo.UpdateResult, error) {
	database := r.Database(r.DatabaseName)
	collection := database.Collection(collectionName)
	filter := bson.D{
		{ID, userID},
		{PASSWORD, password},
		{DELETED, false},
		{ACTIVATED, true},
		{LOCKED, true},
	}
	update := bson.D{
		{"$set", bson.D{
			{LASTEDITED, time.Now()},
			{LOCKED, false},
		}},
	}
	res, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *Repository) dataExists(ctx context.Context, collection *mongo.Collection, fieldName, fieldValue string) bool {
	filter := bson.D{
		{fieldName, fieldValue},
	}
	res := collection.FindOne(ctx, filter)
	user := new(User)
	res.Decode(&user)
	if user.Email != "" {
		return false
	}

	return true
}
