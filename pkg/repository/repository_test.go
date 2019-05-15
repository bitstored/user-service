package repository

import (
	"context"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"
)

type Test struct {
	Name  string
	Field string
	Fails bool
}

const (
	Host           = "mongodb://localhost:27017"
	usernameGood   = "username"
	emailGood      = "email@email.com"
	phoneGood      = "123456789"
	databaseName   = "bitsored"
	collectionName = "users"
	emailGood1     = "email@emailuc.com"
	usernameGood1  = "username1"
	phoneNumber1   = "12345678"
	emailGood2     = "email@emailic.com"
	usernameGood2  = "username2"
	phoneNumber2   = "1234567"
)

func TestNewRepository(t *testing.T) {
	test := Test{"TestCreate", "fieldCreate", false}
	rsp, err := NewRepository(Host, test.Name)
	require.NoError(t, err)
	t.Log(rsp)
}

func TestCreateAccount(t *testing.T) {

	ts := []struct {
		Name    string
		User1   User
		Success bool
	}{
		{"TestCreateAccountSuccess",
			User{
				ID:          primitive.NewObjectID(),
				Email:       emailGood,
				Username:    usernameGood,
				PhoneNumber: phoneGood,
			},
			true,
		},
		{"TestCreateAccountExistsMail",
			User{
				ID:          primitive.NewObjectID(),
				Email:       emailGood,
				Username:    "username1",
				PhoneNumber: "12345678",
			},
			false,
		},
		{"TestCreateAccountExistsUsername",
			User{
				ID:          primitive.NewObjectID(),
				Email:       "email@emailik.com",
				Username:    usernameGood,
				PhoneNumber: "12345678",
			},
			false,
		},
		{"TestCreateAccountExistsPhone",
			User{
				ID:          primitive.NewObjectID(),
				Email:       "emailik@emailik.com",
				Username:    "username23",
				PhoneNumber: phoneGood,
			},
			false,
		},
	}

	for _, tc := range ts {
		tc1 := tc
		t.Run(tc1.Name, func(t *testing.T) {
			r, err := NewRepository(Host, databaseName)
			require.NoError(t, err)
			res, err := r.CreateAccount(context.TODO(), collectionName, tc1.User1)
			if tc1.Success {
				require.NoError(t, err)
				require.NotNil(t, res.InsertedID)

			} else {
				require.Errorf(t, err, "res: %v", res)
				require.Nil(t, res)
			}
		})
	}

	//CLEANUP
	repo, err := NewRepository(Host, databaseName)
	defer repo.Disconnect(context.TODO())

	require.NoError(t, err)
	db := repo.Database(databaseName)
	c := db.Collection(collectionName)
	f := bson.D{
		{USERNAME, usernameGood},
		{EMAIL, emailGood},
		{PHONENUMBER, phoneGood},
	}
	c.DeleteOne(context.TODO(), f)
}

func TestActivateAccount(t *testing.T) {

	ts := []struct {
		Name    string
		User1   User
		Success bool
	}{
		{"TestCreateAccountActivateSuccess",
			User{
				ID:          primitive.NewObjectID(),
				Email:       emailGood,
				Username:    usernameGood,
				PhoneNumber: phoneGood,
				Activated:   false,
			},
			true,
		},
		{"TestCreateAccountNoActivation",
			User{
				ID:          primitive.NewObjectID(),
				Email:       emailGood1,
				Username:    usernameGood1,
				PhoneNumber: phoneNumber1,
				Activated:   false,
			},
			false,
		},
	}

	for _, tc := range ts {
		tc1 := tc
		t.Run(tc1.Name, func(t *testing.T) {
			r, err := NewRepository(Host, databaseName)
			require.NoError(t, err)
			res, err := r.CreateAccount(context.TODO(), collectionName, tc1.User1)
			time.Sleep(1000 * time.Nanosecond)
			require.NoError(t, err)
			require.NotNil(t, res.InsertedID)
			if tc1.Success {
				res1, err := r.ActivateAccount(context.TODO(), collectionName, tc1.User1.ID)
				require.NoErrorf(t, err, "err: %v\nres: %v", err, res1)
				require.NotNil(t, res1)
				time.Sleep(1000 * time.Nanosecond)

			}
			user := r.GetAccount(context.TODO(), collectionName, tc1.User1.ID)
			require.NotNil(t, user)
			if tc1.Success {
				require.True(t, user.Activated)
				require.False(t, user.Deleted)
				require.False(t, user.Locked)
			} else {
				require.False(t, user.Activated)
				require.False(t, user.Deleted)
				require.False(t, user.Locked)
			}
		})
	}

	//CLEANUP
	repo, err := NewRepository(Host, databaseName)
	defer repo.Disconnect(context.TODO())

	require.NoError(t, err)
	db := repo.Database(databaseName)
	c := db.Collection(collectionName)
	f := bson.D{
		{USERNAME, usernameGood},
		{EMAIL, emailGood},
		{PHONENUMBER, phoneGood},
	}
	c.DeleteOne(context.TODO(), f)
	f = bson.D{
		{USERNAME, usernameGood1},
		{EMAIL, emailGood1},
		{PHONENUMBER, phoneNumber1},
	}
	c.DeleteOne(context.TODO(), f)
}

func TestUpdateAccount(t *testing.T) {

	ts := []struct {
		Name    string
		User1   User
		Success bool
	}{
		{"TestUpdateAccountSuccess",
			User{
				ID:          primitive.NewObjectID(),
				Email:       emailGood,
				Username:    usernameGood,
				PhoneNumber: phoneGood,
				Activated:   true,
			},
			true,
		},
		{"TestUpdateAccountNoActivation",
			User{
				ID:          primitive.NewObjectID(),
				Email:       emailGood1,
				Username:    usernameGood1,
				PhoneNumber: phoneNumber1,
				Activated:   false,
			},
			false,
		},
	}

	for _, tc := range ts {
		tc1 := tc
		t.Run(tc1.Name, func(t *testing.T) {
			r, err := NewRepository(Host, databaseName)
			require.NoError(t, err)
			res, err := r.CreateAccount(context.TODO(), collectionName, tc1.User1)
			time.Sleep(1000 * time.Nanosecond)
			require.NoError(t, err)
			require.NotNil(t, res.InsertedID)
			if tc1.Success {
				tc1.User1.LastName = "LasstName"
				tc1.User1.FistName = "FirsstName"

				res1, err := r.UpdateAccount(context.TODO(), collectionName, tc1.User1.ID, tc1.User1)
				require.NoErrorf(t, err, "err: %v\nres: %v", err, res1)
				require.NotNil(t, res1)
				time.Sleep(1000 * time.Nanosecond)
			} else {
				tc1.User1.LastName = "LasstName"
				tc1.User1.FistName = "FirsstName"

				res1, err := r.UpdateAccount(context.TODO(), collectionName, tc1.User1.ID, tc1.User1)
				require.Errorf(t, err, "err: %v\nres: %v", err, res1)
				require.Nil(t, res1)
				time.Sleep(1000 * time.Nanosecond)
			}
			user := r.GetAccount(context.TODO(), collectionName, tc1.User1.ID)
			require.NotNil(t, user)
			if tc1.Success {
				require.True(t, user.Activated)
				require.False(t, user.Deleted)
				require.False(t, user.Locked)
				require.Equal(t, tc1.User1.FistName, user.FistName)
				require.Equal(t, tc1.User1.LastName, user.LastName)
			} else {
				require.False(t, user.Activated)
				require.False(t, user.Deleted)
				require.False(t, user.Locked)
				require.Empty(t, user.FistName)
				require.Empty(t, user.LastName)
			}
		})
	}

	//CLEANUP
	repo, err := NewRepository(Host, databaseName)
	defer repo.Disconnect(context.TODO())

	require.NoError(t, err)
	db := repo.Database(databaseName)
	c := db.Collection(collectionName)
	f := bson.D{
		{USERNAME, usernameGood},
		{EMAIL, emailGood},
		{PHONENUMBER, phoneGood},
	}
	c.DeleteOne(context.TODO(), f)
	f = bson.D{
		{USERNAME, usernameGood1},
		{EMAIL, emailGood1},
		{PHONENUMBER, phoneNumber1},
	}
	c.DeleteOne(context.TODO(), f)
}

func TestDeleteAccount(t *testing.T) {

	ts := []struct {
		Name    string
		User1   User
		Success bool
	}{
		{"TestDeleteAccountInvalidPass",
			User{
				ID:          primitive.NewObjectID(),
				Email:       emailGood2,
				Username:    usernameGood2,
				PhoneNumber: phoneNumber2,
				Activated:   true,
				Password:    "pass1",
			},
			false,
		},
		{"TestDeleteAccountSuccess",
			User{
				ID:          primitive.NewObjectID(),
				Email:       emailGood,
				Username:    usernameGood,
				PhoneNumber: phoneGood,
				Activated:   true,
				Password:    "pass",
			},
			true,
		},
		{"TestDeleteAccountNoActivation",
			User{
				ID:          primitive.NewObjectID(),
				Email:       emailGood1,
				Username:    usernameGood1,
				PhoneNumber: phoneNumber1,
				Activated:   false,
				Password:    "pass",
			},
			false,
		},
	}

	for _, tc := range ts {
		tc1 := tc
		t.Run(tc1.Name, func(t *testing.T) {
			r, err := NewRepository(Host, databaseName)
			require.NoError(t, err)
			res, err := r.CreateAccount(context.TODO(), collectionName, tc1.User1)
			time.Sleep(100 * time.Nanosecond)
			require.NoError(t, err)
			require.NotNil(t, res.InsertedID)
			if tc1.Success {
				res1, err := r.DeleteAccount(context.TODO(), collectionName, tc1.User1.ID, "pass")
				require.NoErrorf(t, err, "err: %v\nres: %v", err, res1)
				require.NotNil(t, res1)
				time.Sleep(100 * time.Nanosecond)
			} else {
				res1, err := r.DeleteAccount(context.TODO(), collectionName, tc1.User1.ID, "pass")
				require.Errorf(t, err, "err: %v\nres: %v", err, res1)
				require.Nil(t, res1)
				time.Sleep(100 * time.Nanosecond)
			}
			user := r.GetAccount(context.TODO(), collectionName, tc1.User1.ID)
			require.NotNil(t, user)
			if tc1.Success {
				require.True(t, user.Activated)
				require.True(t, user.Deleted)
				require.False(t, user.Locked)
			} else {
				require.Equal(t, tc1.User1.Activated, user.Activated)
				require.False(t, user.Deleted)
				require.False(t, user.Locked)
			}
		})
	}

	//CLEANUP
	repo, err := NewRepository(Host, databaseName)
	defer repo.Disconnect(context.TODO())

	require.NoError(t, err)
	db := repo.Database(databaseName)
	c := db.Collection(collectionName)
	f := bson.D{
		{USERNAME, usernameGood},
		{EMAIL, emailGood},
		{PHONENUMBER, phoneGood},
	}
	c.DeleteOne(context.TODO(), f)
	f = bson.D{
		{USERNAME, usernameGood1},
		{EMAIL, emailGood1},
		{PHONENUMBER, phoneNumber1},
	}
	c.DeleteOne(context.TODO(), f)
	f = bson.D{
		{USERNAME, usernameGood2},
		{EMAIL, emailGood2},
		{PHONENUMBER, phoneNumber2},
	}
	c.DeleteOne(context.TODO(), f)
}
