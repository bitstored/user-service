package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/bitstored/user-service/pb"
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
		Repo:        repo,
		Activations: make(map[string]Activation),
		Sessions:    make(map[string]Session),
	}
}

func (s *Service) CreateAccount(ctx context.Context, fName, lName string, bDay time.Time, email, uname, pass, pNumber string, photo string) (string, error) {

	_, err1 := validator.Password(pass)
	if err1 != nil {
		return "", err1.Error()
	}

	ok := validator.Email(email)
	if !ok {
		return "", fmt.Errorf("Email is invalid")
	}
	// Encrypt pass
	salt := newSalt()
	hash, err := encryptPassword(pass, salt)

	if err != nil {
		return "", err
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
		return "", err
	}

	if res.InsertedID == nil {
		return "", fmt.Errorf("unable to create user")
	}

	// Generate activation token
	token := make([]byte, ACTIVATION_TOKEN_LEN)
	rand.Read(token)
	s.Activations[string(token)] = Activation{email, string(token), time.Now().Add(EXPIRES_AFTER_MINUTES_PERIOD * time.Minute)}

	// send activation mail
	// smtpServer := NewSmtpServer(SMTP_HOST, SMTP_PORT)
	// messages := []Message{{to: email, sender: SMTP_ADMIN_USERNAME, password: SMTP_ADMIN_PASSWORD, subject: "Activation mail", link: HOST + "/activate/" + string(token)}}

	// err = smtpServer.sendMail(ctx, messages)

	return user.ID.String(), nil
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

func (s *Service) UpdateAccount(ctx context.Context, token, password, firstname, lastname string, photo string) error {
	session, ok := s.Sessions[token]
	fmt.Printf("Sessions %v\n\n Session %v\n\n\n", s.Sessions, session)
	if !ok {
		return fmt.Errorf("Session token is invalid")
	}
	// ok = s.validateToken(ctx, token, session.ID, session.FirstName, session.LastName)
	// if !ok {
	// 	return nil, fmt.Errorf("Session token is invalid")
	// }
	u := s.Repo.GetAccount(ctx, USER_COLLECTION_NAME, session.ID)
	if u == nil {
		return fmt.Errorf("Session token is invalid")
	}
	salt := u.Salt
	password, err := encryptPassword(password, salt)
	if err != nil {
		return err
	}
	user := repository.User{
		ID:        session.ID,
		FirstName: firstname,
		LastName:  lastname,
		Password:  password,
		Photo:     photo,
	}
	res, err := s.Repo.UpdateAccount(ctx, USER_COLLECTION_NAME, session.ID, user)
	if res.MatchedCount == 0 {
		return fmt.Errorf("Unable to update user")
	}
	return err
}

func (s *Service) DeleteAccount(ctx context.Context, token, password string) (bool, error) {
	session, ok := s.Sessions[token]
	if !ok {
		return false, fmt.Errorf("Session token is invalid")
	}
	uid := session.ID

	u := s.Repo.GetAccount(ctx, USER_COLLECTION_NAME, session.ID)
	if u == nil {
		return false, fmt.Errorf("Session token is invalid")
	}
	salt := u.Salt
	password, err := encryptPassword(password, salt)
	if err != nil {
		return false, err
	}
	res, err := s.Repo.DeleteAccount(ctx, USER_COLLECTION_NAME, uid, password)
	if err != nil || res.ModifiedCount == 0 {
		return false, err
	}
	return true, nil
}

func (s *Service) GetAccount(ctx context.Context, token string) (*pb.User, error) {
	session, ok := s.Sessions[token]
	fmt.Printf("Sessions %v\n\n Session %v\n\n\n", s.Sessions, session)
	if !ok {
		return nil, fmt.Errorf("Session token is invalid")
	}
	// ok = s.validateToken(ctx, token, session.ID, session.FirstName, session.LastName)
	// if !ok {
	// 	return nil, fmt.Errorf("Session token is invalid")
	// }
	fmt.Printf("%v\n\n", session.ID)
	u := s.Repo.GetAccount(ctx, USER_COLLECTION_NAME, session.ID)
	if u == nil {
		return nil, fmt.Errorf("User not found")
	}
	user := &pb.User{
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		Email:       u.Email,
		Username:    u.Username,
		PhoneNumber: u.PhoneNumber,
		Photo:       u.Photo,
		Birthday:    u.Birthday.String(),
		LastLogin:   u.LastLogin.String(),
		LastEdited:  u.LastEdited.String(),
		Created:     session.ID.String(),
	}
	return user, nil
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

	token, err := s.createToken(ctx, *user)
	if err != nil {
		return "", err
	}
	s.Sessions[token] = Session{ID: user.ID, FirstName: user.FirstName, LastName: user.LastName, ExpiresAt: time.Now().Add(time.Hour * EXPIRES_AFTER_HOURS_PERIOD)}

	return token, nil
}

func (s *Service) Logout(ctx context.Context, token string) (bool, error) {
	if _, ok := s.Sessions[token]; ok {
		fmt.Println("Loging out")
		delete(s.Sessions, token)
		return true, nil
	}
	fmt.Println("No token")
	return false, fmt.Errorf("Loging out, invalid token")
}

func (s *Service) ResetPassword(ctx context.Context, token, password, newPassword string) (bool, error) {
	session, ok := s.Sessions[token]
	if !ok {
		return false, fmt.Errorf("Session token invalid")
	}
	_, err1 := validator.Password(newPassword)
	if err1 != nil {
		return false, err1.Error()
	}

	u := s.Repo.GetAccount(ctx, USER_COLLECTION_NAME, session.ID)
	if u == nil {
		return false, fmt.Errorf("User not found")
	}
	newPassword, err := encryptPassword(newPassword, u.Salt)
	if err != nil {
		return false, err
	}
	res, err := s.Repo.ResetPassword(ctx, USER_COLLECTION_NAME, session.ID.String(), u.Password, newPassword)
	if err != nil {
		return false, err
	}
	if res.ModifiedCount != 1 {
		return false, fmt.Errorf("User not found")
	}
	return true, nil
}
func (s *Service) LockAccount(ctx context.Context, token, userID string) (bool, error) {

	session, ok := s.Sessions[token]
	if !ok {
		return false, fmt.Errorf("Session token invalid")
	}

	u := s.Repo.GetAccount(ctx, USER_COLLECTION_NAME, session.ID)
	if u == nil {
		return false, fmt.Errorf("User not found")
	}

	if session.ID.String() == userID || u.IsAdmin {
		res, err := s.Repo.LockAccount(ctx, USER_COLLECTION_NAME, userID, "")
		if err != nil {
			return false, err
		}
		if res.ModifiedCount != 1 {
			return false, fmt.Errorf("User not found")
		}
	}

	return true, nil
}
func (s *Service) RequestUnlockAccount(ctx context.Context) (bool, error) {

	return false, nil
}
func (s *Service) UnlockAccount(ctx context.Context, token, userID string) (bool, error) {

	session, ok := s.Sessions[token]
	if !ok {
		return false, fmt.Errorf("Session token invalid")
	}

	u := s.Repo.GetAccount(ctx, USER_COLLECTION_NAME, session.ID)
	if u == nil {
		return false, fmt.Errorf("User not found")
	}

	if session.ID.String() == userID || u.IsAdmin {
		res, err := s.Repo.UnlockAccount(ctx, USER_COLLECTION_NAME, userID, "")
		if err != nil {
			return false, err
		}
		if res.ModifiedCount != 1 {
			return false, fmt.Errorf("User not found")
		}
	}

	return true, nil
}

func (s *Service) ListUsers(ctx context.Context, token string) ([]*pb.User, error) {

	session, ok := s.Sessions[token]
	if !ok {
		return nil, fmt.Errorf("Session token invalid")
	}

	u := s.Repo.GetAccount(ctx, USER_COLLECTION_NAME, session.ID)
	if u == nil {
		return nil, fmt.Errorf("User not found")
	}
	if !u.IsAdmin {
		return nil, fmt.Errorf("Access denied")
	}

	users, err := s.Repo.ListUsers(ctx, USER_COLLECTION_NAME)
	if err != nil {
		return nil, fmt.Errorf("Access denied")
	}
	pbUsers := make([]*pb.User, 0)
	for _, u := range users {
		U := new(pb.User)
		U.AvailableMb = u.AvailableMb
		U.Birthday = u.Birthday.String()
		U.Email = u.Email
		U.FirstName = u.FirstName
		U.LastName = u.LastName
		U.PhoneNumber = u.PhoneNumber
		U.Photo = u.Photo
		U.Username = u.Username
		U.IsActivated = u.Activated
		U.IsAdmin = u.IsAdmin
		U.IsLocked = u.Locked
		U.LastEdited = u.LastEdited.String()
		U.LastLogin = u.LastLogin.String()
		U.Created = u.CreatedAt.String()
		pbUsers = append(pbUsers, U)
	}
	return pbUsers, nil
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

func (s *Service) createToken(ctx context.Context, user repository.User) (string, error) {

	conn, err := grpc.Dial(AUTH_GRPC_PORT, grpc.WithInsecure())
	defer conn.Close()

	if err != nil {
		return "", err
	}
	client := auth.NewAuthServiceClient(conn)

	req := &auth.GenerateJWTRequest{UserId: user.ID.String(), FirstName: user.FirstName, Lastname: user.LastName, IsAdmin: user.IsAdmin}
	fmt.Printf("rerq %v\n\n", req)
	rsp, err := client.GenerateJWT(ctx, req)
	if err != nil {
		return "", err
	}
	//Add token to user tokens
	return rsp.Token, nil
}
