package service

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/smtp"
)

type Message struct {
	to       string
	sender   string
	password string
	subject  string
	body     string
	link     string
}

type SmtpServer struct {
	host string
	port string
}

func (s *SmtpServer) ServerName() string {
	return s.host + ":" + s.port
}

func NewSmtpServer(host, port string) *SmtpServer {
	return &SmtpServer{host, port}
}

func (m *Message) BuildMessage() string {
	message := ""
	message += "content-type: text/html \r\n"
	message += fmt.Sprintf("From: %s\r\n", m.sender)
	message += fmt.Sprintf("To: %s\r\n", m.to)

	message += fmt.Sprintf("Subject: %s\r\n", m.subject)
	message += "\r\n <h3>" + m.body + "</h3>"
	message += "\r\n <a style=\"color:#0000FF\" href=" + m.link + "> Confirmation Link </a>"
	message += "\r\n <h4 style=\"color:#FF0000\">Please don't reply to this mail </h4>"
	return message
}

func (s *SmtpServer) sendMail(ctx context.Context, messages []Message) error {

	// Set the sender and recipient.
	for _, m := range messages {
		auth := smtp.PlainAuth("", m.sender, m.password, s.host)
		mess := m.BuildMessage()

		tlsconfig := &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         s.host,
		}
		// Send the email body.
		conn, err := tls.Dial("tcp", s.ServerName(), tlsconfig)
		if err != nil {
			return err
		}

		client, err := smtp.NewClient(conn, s.host)
		if err != nil {
			return err
		}

		// step 1: Use Auth
		if err = client.Auth(auth); err != nil {
			return err
		}

		// step 2: add all from and to
		if err = client.Mail(m.sender); err != nil {
			return err
		}

		if err = client.Rcpt(m.to); err != nil {
			return err
		}

		// Data
		w, err := client.Data()
		if err != nil {
			return err
		}

		_, err = w.Write([]byte(mess))
		if err != nil {
			return err
		}

		err = w.Close()
		if err != nil {
			return err
		}

		client.Quit()

		fmt.Println("Mail sent successfully")
	}
	return nil
}
