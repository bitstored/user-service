package server

import (
	"context"
	"flag"
	"google.golang.org/grpc"
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/bitstored/user-service/pb"
	"github.com/bitstored/user-service/pkg/repository"
	"github.com/bitstored/user-service/pkg/service"
	"github.com/stretchr/testify/require"
)

const (
	repoHost = "mongodb://localhost:27017"
)

var (
	grpcAddr = flag.String("grpc", "localhost:4008", "gRPC API address")
)

func TestServer_Login(t *testing.T) {
	username := "testlogin"
	password := "Alorap1!"
	//email := "test@test.com"
	repo, err := repository.NewRepository(repoHost, t.Name())
	serv := service.NewService(repo)
	require.NoError(t, err)
	gRPCListener, err := net.Listen("tcp", *grpcAddr)
	devServer := NewServer(serv)

	// Register standard server metrics and customized metrics to registry.

	gRPCServer := grpc.NewServer()

	pb.RegisterAccountServer(gRPCServer, devServer)
	go func() {
		if err := gRPCServer.Serve(gRPCListener); err != nil {
			require.NoErrorf(t, err, "Failed to serve gRPC: %s", err)
		}
	}()
	ctx := context.TODO()
	conn, err := grpc.DialContext(ctx, *grpcAddr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	go func() {
		<-ctx.Done()
		if err := conn.Close(); err != nil {
			require.NoErrorf(t, err, "Failed to close a client connection to the gRPC server: %v", err)
		}
	}()

	type fields struct {
		Service *service.Service
	}
	type args struct {
		ctx context.Context
		in  *pb.LoginRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.LoginResponse
		wantErr bool
	}{
		{
			name: t.Name() + "Success",
			fields: fields{
				Service: serv,
			},
			args: args{
				ctx: context.TODO(),
				in: &pb.LoginRequest{
					Username: username,
					Password: password,
				},
			},
			want:    &pb.LoginResponse{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			conn1, err := grpc.DialContext(context.TODO(), "localhost:4008", grpc.WithInsecure())
			require.NoErrorf(t, err, "failed to dial bufnet: %v", err)
			defer conn1.Close()
			client := pb.NewAccountClient(conn1)

			got, err := client.Login(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Server.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Server.Login() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewServer(t *testing.T) {
	type args struct {
		s *service.Service
	}
	tests := []struct {
		name string
		args args
		want *Server
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewServer(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_CreateAccount(t *testing.T) {
	type fields struct {
		Service *service.Service
	}
	type args struct {
		ctx context.Context
		in  *pb.CreateAccountRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.CreateAccountResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				Service: tt.fields.Service,
			}
			got, err := s.CreateAccount(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Server.CreateAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Server.CreateAccount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_ResendActivationMail(t *testing.T) {
	type fields struct {
		Service *service.Service
	}
	type args struct {
		ctx context.Context
		in  *pb.ResendActivationMailRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.ResendActivationMailResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				Service: tt.fields.Service,
			}
			got, err := s.ResendActivationMail(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Server.ResendActivationMail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Server.ResendActivationMail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_ActivateAccount(t *testing.T) {
	type fields struct {
		Service *service.Service
	}
	type args struct {
		ctx context.Context
		in  *pb.ActivateAccountRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.ActivateAccountResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				Service: tt.fields.Service,
			}
			got, err := s.ActivateAccount(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Server.ActivateAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Server.ActivateAccount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_UpdateAccount(t *testing.T) {
	type fields struct {
		Service *service.Service
	}
	type args struct {
		ctx context.Context
		in  *pb.UpdateAccountRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.UpdateAccountResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				Service: tt.fields.Service,
			}
			got, err := s.UpdateAccount(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Server.UpdateAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Server.UpdateAccount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_DeleteAccount(t *testing.T) {
	type fields struct {
		Service *service.Service
	}
	type args struct {
		ctx context.Context
		in  *pb.DeleteAccountRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.DeleteAccountResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				Service: tt.fields.Service,
			}
			got, err := s.DeleteAccount(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Server.DeleteAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Server.DeleteAccount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_GetAccount(t *testing.T) {
	type fields struct {
		Service *service.Service
	}
	type args struct {
		ctx context.Context
		in  *pb.GetAccountRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.GetAccountResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				Service: tt.fields.Service,
			}
			got, err := s.GetAccount(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Server.GetAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Server.GetAccount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_Logout(t *testing.T) {
	type fields struct {
		Service *service.Service
	}
	type args struct {
		ctx context.Context
		in  *pb.LogoutRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.LogoutResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				Service: tt.fields.Service,
			}
			got, err := s.Logout(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Server.Logout() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Server.Logout() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_ResetPassword(t *testing.T) {
	type fields struct {
		Service *service.Service
	}
	type args struct {
		ctx context.Context
		in  *pb.ResetPasswordRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.ResetPasswordResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				Service: tt.fields.Service,
			}
			got, err := s.ResetPassword(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Server.ResetPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Server.ResetPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_LockAccount(t *testing.T) {
	type fields struct {
		Service *service.Service
	}
	type args struct {
		ctx context.Context
		in  *pb.LockAccountRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.LockAccountResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				Service: tt.fields.Service,
			}
			got, err := s.LockAccount(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Server.LockAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Server.LockAccount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_RequestUnlockAccount(t *testing.T) {
	type fields struct {
		Service *service.Service
	}
	type args struct {
		ctx context.Context
		in  *pb.RequestUnlockAccountRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.RequestUnlockAccountResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				Service: tt.fields.Service,
			}
			got, err := s.RequestUnlockAccount(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Server.RequestUnlockAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Server.RequestUnlockAccount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_UnlockAccount(t *testing.T) {
	type fields struct {
		Service *service.Service
	}
	type args struct {
		ctx context.Context
		in  *pb.UnlockAccountRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.UnlockAccountResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				Service: tt.fields.Service,
			}
			got, err := s.UnlockAccount(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Server.UnlockAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Server.UnlockAccount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_ListUsers(t *testing.T) {
	type fields struct {
		Service *service.Service
	}
	type args struct {
		ctx context.Context
		in  *pb.ListUsersRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.ListUsersResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				Service: tt.fields.Service,
			}
			got, err := s.ListUsers(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Server.ListUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Server.ListUsers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateUser(t *testing.T) {
	type args struct {
		user pb.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateUser(tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("validateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_parseDate(t *testing.T) {
	type args struct {
		date string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseDate(tt.args.date)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseDate() = %v, want %v", got, tt.want)
			}
		})
	}
}
