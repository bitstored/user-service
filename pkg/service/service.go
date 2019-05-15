package service

import (
	"context"
	"fmt"
	"time"

	"github.com/bitstored/auth-service/pkg/validator"
	"github.com/bitstored/user-service/pkg/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	USER_COLLECTION_NAME      = "users"
	ACCOUNT_COLLECTION_NAME   = "accounts"
	PAYEMENTS_COLLECTION_NAME = "payements"
)

type Service struct {
	Repo        *repository.Repository
	Activations map[primitive.ObjectID]Activation
	Sessions    map[primitive.ObjectID]Session
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Repo: repo,
	}
}

func (s *Service) CreateAccount(ctx context.Context, fName, lName string, bDay time.Time, email, uname, pass, pNumber string, photo []byte) error {

	_, err1 := validator.Password(pass)
	if err1 != nil {
		return err1.Error()
	}

	ok := validator.Email(email)
	if !ok {
		return fmt.Errorf("Email is invalid")
	}

	user := repository.User{
		ID:          primitive.NewObjectID(),
		FistName:    fName,
		LastName:    lName,
		Birthday:    bDay,
		Email:       email,
		Username:    uname,
		Password:    pass,
		PhoneNumber: pNumber,
		Photo:       photo,
	}

	res, err := s.Repo.CreateAccount(ctx, USER_COLLECTION_NAME, user)

	if err != nil {
		return err
	}

	if res.InsertedID == nil {
		return fmt.Errorf("unable to create user")
	}

	return nil
}

func (s *Service) ResendActivationMail(ctx context.Context) (error, error) {
	return nil, nil
}
func (s *Service) ActivateAccount(ctx context.Context) (error, error) {
	return nil, nil
}
func (s *Service) UpdateAccount(ctx context.Context) (error, error) {
	return nil, nil
}
func (s *Service) DeleteAccount(ctx context.Context) (error, error) {
	return nil, nil
}
func (s *Service) GetAccount(ctx context.Context) (error, error) {
	return nil, nil
}
func (s *Service) Login(ctx context.Context) (error, error) {
	return nil, nil
}
func (s *Service) Logout(ctx context.Context) (error, error) {
	return nil, nil
}
func (s *Service) ResetPassword(ctx context.Context) (error, error) {
	return nil, nil
}
func (s *Service) LockAccount(ctx context.Context) (error, error) {
	return nil, nil
}
func (s *Service) RequestUnlockAccount(ctx context.Context) (error, error) {
	return nil, nil
}
func (s *Service) UnlockAccount(ctx context.Context) (error, error) {
	return nil, nil
}
