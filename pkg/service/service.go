package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	auth "github.com/bitstored/auth-service/pb"
	"github.com/bitstored/auth-service/pkg/validator"
	"github.com/bitstored/user-service/pkg/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
)

const (
	USER_COLLECTION_NAME         = "users"
	ACCOUNT_COLLECTION_NAME      = "accounts"
	PAYEMENTS_COLLECTION_NAME    = "payements"
	SMTP_HOST                    = "smtp.gmail.com"
	SMTP_PORT                    = "465"
	SMTP_ADMIN_USERNAME          = "hennessyparadise@gmail.com"
	SMTP_ADMIN_PASSWORD          = "dianuscabejan1996"
	EXPIRES_AFTER_MINUTES_PERIOD = 60
	EXPIRES_AFTER_HOURS_PERIOD   = 7 * 24
	AUTH_GRPC_PORT               = "localhost:4002"
	ACTIVATION_TOKEN_LEN         = 32
	HOST                         = "localhost:5008"
)

type Service struct {
	Repo        *repository.Repository
	Activations map[string]Activation
	Sessions    map[string]Session
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
	// Encrypt pass
	salt := newSalt()
	hash, err := encryptPassword(pass, salt)

	if err != nil {
		return err
	}
	user := repository.User{
		ID:          primitive.NewObjectID(),
		FirstName:   fName,
		LastName:    lName,
		Birthday:    bDay,
		Email:       email,
		Username:    uname,
		Password:    hash,
		Salt:        salt,
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

	// Generate activation token
	token := make([]byte, ACTIVATION_TOKEN_LEN)
	rand.Read(token)
	s.Activations[string(token)] = Activation{email, string(token), time.Now().Add(EXPIRES_AFTER_MINUTES_PERIOD * time.Minute)}

	// send activation mail
	smtpServer := NewSmtpServer(SMTP_HOST, SMTP_PORT)
	messages := []Message{{to: email, sender: SMTP_ADMIN_USERNAME, password: SMTP_ADMIN_PASSWORD, subject: "Activation mail", link: HOST + "/user/api/v1/account/activate/" + string(token)}}

	err = smtpServer.sendMail(ctx, messages)

	return err
}

func (s *Service) ResendActivationMail(ctx context.Context, email string) error {
	// Generate activation token
	token := make([]byte, ACTIVATION_TOKEN_LEN)
	rand.Read(token)
	s.Activations[string(token)] = Activation{email, string(token), time.Now().Add(EXPIRES_AFTER_MINUTES_PERIOD * time.Minute)}

	// send activation mail
	smtpServer := NewSmtpServer(SMTP_HOST, SMTP_PORT)
	messages := []Message{{to: email, sender: SMTP_ADMIN_USERNAME, password: SMTP_ADMIN_PASSWORD, subject: "Activation mail", link: HOST + "/user/api/v1/account/activate/" + string(token)}}

	err := smtpServer.sendMail(ctx, messages)
	return err
}

func (s *Service) ActivateAccount(ctx context.Context, token string) error {
	if activation, ok := s.Activations[token]; ok {
		if activation.ExpiresAt.Before(time.Now()) {
			s.ResendActivationMail(ctx, activation.Email)
			return fmt.Errorf("Token Expired, activation mail resend")
		}
		s.Repo.ActivateAccount(ctx, USER_COLLECTION_NAME, activation.Email)
	}
	return fmt.Errorf("Unable to activate account, token not found")
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

func (s *Service) Login(ctx context.Context, username, password string) (string, error) {
	if username == "" {
		return "", fmt.Errorf("empty username")
	}

	user := s.Repo.Login(ctx, USER_COLLECTION_NAME, username, password)
	// Encrypt pass
	if user == nil {
		return "", fmt.Errorf("User not found")
	}
	salt := user.Salt

	hash, err := encryptPassword(password, salt)

	if err != nil {
		return "", err
	}

	if hash != user.Password {
		return "", fmt.Errorf("invalid password")
	}
	//TODO add attempts

	token := s.createToken(ctx, *user)

	s.Sessions[token] = Session{ID: user.ID, FirstName: user.FirstName, LastName: user.LastName, ExpiresAt: time.Now().Add(time.Hour * EXPIRES_AFTER_HOURS_PERIOD)}

	return token, nil
}

func (s *Service) Logout(ctx context.Context, token string) (bool, error) {
	if session, ok := s.Sessions[token]; ok {
		valid := s.validateToken(ctx, token, session.ID, session.FirstName, session.LastName)
		if valid {
			delete(s.Sessions, token)
			return true, nil
		}
	}
	return false, fmt.Errorf("Unable to logout")
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

func (s *Service) addActivation(email string, code string) {

	expiresAt := time.Now().Add(time.Minute * EXPIRES_AFTER_MINUTES_PERIOD)

	a := Activation{Email: email, Token: code, ExpiresAt: expiresAt}

	s.Activations[email] = a
}

func (s *Service) canActivate(email string, code string) bool {

	now := time.Now()
	a := s.Activations[email]
	if now.After(a.ExpiresAt) {
		delete(s.Activations, email)
		return false
	}

	if code != a.Token {
		return false
	}

	delete(s.Activations, email)
	return true
}

func (s *Service) validateToken(ctx context.Context, token string, uid primitive.ObjectID, firstname, lastname string) bool {
	conn, err := grpc.Dial(AUTH_GRPC_PORT, grpc.WithInsecure())
	defer conn.Close()

	if err != nil {
		return false
	}
	client := auth.NewAuthServiceClient(conn)

	req := &auth.ValidateJWTRequest{Token: token, UserId: uid.String(), FirstName: firstname, Lastname: lastname}

	rsp, err := client.ValidateJWT(ctx, req)
	if err != nil {
		return false
	}
	return rsp.GetIsValid()
}

func (s *Service) createToken(ctx context.Context, user repository.User) string {

	conn, err := grpc.Dial(AUTH_GRPC_PORT, grpc.WithInsecure())
	defer conn.Close()

	if err != nil {
		return ""
	}
	client := auth.NewAuthServiceClient(conn)

	req := &auth.GenerateJWTRequest{UserId: user.ID.String(), FirstName: user.FirstName, Lastname: user.LastName, IsAdmin: user.IsAdmin}

	rsp, err := client.GenerateJWT(ctx, req)
	if err != nil {
		return ""
	}
	//Add token to user tokens
	return rsp.Token
}
