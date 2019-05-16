package service

import (
	"context"
	"crypto/rand"
	crypto "github.com/bitstored/crypto-service/pb"
	"google.golang.org/grpc"
)

const (
	gRPCPortCrypto = "localhost:4004"
	iterationCount = 1000
	saltLen        = 32
)

func encryptPassword(password string, salt []byte) (string, error) {

	conn, err := grpc.Dial(gRPCPortCrypto, grpc.WithInsecure())
	if err != nil {
		return "", err
	}

	client := crypto.NewCryptoClient(conn)
	req := crypto.EncryptPasswordRequest{Password: []byte(password), Salt: salt, IterationCount: iterationCount}
	rsp, err := client.EncryptPassword(context.TODO(), &req)

	if err != nil {
		return "", err
	}

	return string(rsp.GetPassword()), nil
}

func newSalt() []byte {
	b := make([]byte, saltLen)

	_, err := rand.Read(b)

	if err != nil {
		return nil
	}

	return b
}
