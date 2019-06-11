package server

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"

	file_pb "github.com/bitstored/file-service/pb"
	"github.com/bitstored/user-service/pb"
	"github.com/bitstored/user-service/pkg/service"
)

const (
	FILE_GRPC_PORT = "localhost:4005"
)

type Server struct {
	Service *service.Service
}

func NewServer(s *service.Service) *Server {
	return &Server{s}
}

func (s *Server) CreateAccount(ctx context.Context, in *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {

	user := *in.GetUser()

	fmt.Printf("user %v \n\n\n\n\n", in.GetUser())
	err := validateUser(user)
	if err != nil {
		return nil, err
	}

	date, err := parseDate(user.GetBirthday())
	if err != nil {
		return nil, err
	}

	uid, err := s.Service.CreateAccount(ctx, user.GetFirstName(), user.GetLastName(), date, user.GetEmail(), user.GetUsername(), user.GetPassword(), user.GetPhoneNumber(), user.GetPhoto())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	conn, err := grpc.Dial(FILE_GRPC_PORT, grpc.WithInsecure())
	defer conn.Close()

	if err != nil {
		return nil, status.Error(codes.Internal, "Unable to create Drive")
	}
	client := file_pb.NewFileManagementClient(conn)

	client.CreateDrive(ctx, &file_pb.CreateDriveRequest{UserId: uid})
	return &pb.CreateAccountResponse{UserId: uid}, nil
}

func (s *Server) ResendActivationMail(ctx context.Context, in *pb.ResendActivationMailRequest) (*pb.ResendActivationMailResponse, error) {

	err := s.Service.ResendActivationMail(ctx, in.GetEmail())
	if err != nil {
		return nil, err
	}
	return &pb.ResendActivationMailResponse{}, nil
}

func (s *Server) ActivateAccount(ctx context.Context, in *pb.ActivateAccountRequest) (*pb.ActivateAccountResponse, error) {
	err := s.Service.ActivateAccount(ctx, in.GetActivationToken())

	if err != nil {
		return nil, err
	}

	return &pb.ActivateAccountResponse{}, nil
}

func (s *Server) UpdateAccount(ctx context.Context, in *pb.UpdateAccountRequest) (*pb.UpdateAccountResponse, error) {
	return nil, nil
}

func (s *Server) DeleteAccount(ctx context.Context, in *pb.DeleteAccountRequest) (*pb.DeleteAccountResponse, error) {
	return nil, nil
}

func (s *Server) GetAccount(ctx context.Context, in *pb.GetAccountRequest) (*pb.GetAccountResponse, error) {
	token := in.GetId()
	if token == "" {
		return nil, status.Error(codes.InvalidArgument, "User not found")
	}
	user, err := s.Service.GetAccount(ctx, token)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.GetAccountResponse{User: user}, nil
}

func (s *Server) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	token, err := s.Service.Login(ctx, in.GetUsername(), in.GetPassword())
	if err != nil {
		return nil, err
	}
	fmt.Printf("Login %s %s, token %v\n", in.Username, in.Password, token)
	return &pb.LoginResponse{SessionToken: token}, nil
}

func (s *Server) Logout(ctx context.Context, in *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	_, err := s.Service.Logout(ctx, in.GetSessionToken())
	if err != nil {
		return nil, err
	}
	return &pb.LogoutResponse{}, nil
}

func (s *Server) ResetPassword(ctx context.Context, in *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
	_, err := s.Service.ResetPassword(ctx, in.GetSessionToken(), in.GetOldPassword(), in.GetNewPassword())
	if err != nil {
		return nil, err
	}
	return &pb.ResetPasswordResponse{}, nil
}

func (s *Server) LockAccount(ctx context.Context, in *pb.LockAccountRequest) (*pb.LockAccountResponse, error) {
	_, err := s.Service.LockAccount(ctx, in.GetSessionToken(), in.GetUserId())
	if err != nil {
		return nil, err
	}
	return &pb.LockAccountResponse{}, nil
}

func (s *Server) RequestUnlockAccount(ctx context.Context, in *pb.RequestUnlockAccountRequest) (*pb.RequestUnlockAccountResponse, error) {
	return nil, nil
}

func (s *Server) UnlockAccount(ctx context.Context, in *pb.UnlockAccountRequest) (*pb.UnlockAccountResponse, error) {
	_, err := s.Service.UnlockAccount(ctx, in.GetSessionToken(), in.GetUserId())
	if err != nil {
		return nil, err
	}
	return &pb.UnlockAccountResponse{}, nil
}

func (s *Server) ListUsers(ctx context.Context, in *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	users, err := s.Service.ListUsers(ctx, in.GetSessionToken())
	if err != nil {
		return nil, err
	}
	return &pb.ListUsersResponse{Users: users}, nil
}

func validateUser(user pb.User) error {
	if user.GetUsername() == "" {
		return status.Error(codes.InvalidArgument, "username is empty")
	}
	if user.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is empty")
	}
	if user.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is empty")
	}
	if user.GetFirstName() == "" {
		return status.Error(codes.InvalidArgument, "firstname is empty")
	}
	if user.GetLastName() == "" {
		return status.Error(codes.InvalidArgument, "lastname is empty")
	}
	return nil
}

func parseDate(date string) (time.Time, error) {
	layout := "2006-01-02T15:04:05.000Z"
	t, err := time.Parse(layout, date)

	return t, err
}
