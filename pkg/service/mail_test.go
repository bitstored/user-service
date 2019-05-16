package service

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSendMail(t *testing.T) {
	messages := []Message{
		{
			sender:   "hennessyparadise2@gmail.com",
			password: "dianuscabejan1996",
			to:       "bejan.diana.andrei@gmail.com",
			body:     "hey there",
			link:     "https://hackernoon.com/golang-sendmail-sending-mail-through-net-smtp-package-5cadbe2670e0",
		},
	}

	server := NewSmtpServer("smtp.gmail.com", "465")
	err := server.sendMail(context.TODO(), messages)
	require.NoError(t, err)
}
