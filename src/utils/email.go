package utils

import (
	"errors"
	"github.com/aja-video/contra/src/configuration"
	"net/smtp"
	"strconv"
)

// SendEmail to user
// - Consider using a library.
func SendEmail(c *configuration.Config, subject, message string) error {

	auth := buildAuth(c.SMTPUser, c.SMTPPass)

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	err := smtp.SendMail(
		c.SMTPHost+":"+strconv.Itoa(c.SMTPPort),
		auth,
		c.EmailFrom,
		[]string{c.EmailTo},
		[]byte(
			"To: <"+c.EmailTo+">\r\n"+
				"From: ContraMail <"+c.EmailFrom+">\r\n"+
				"Subject: "+subject+"\r\n"+
				"\r\n"+
				message,
		),
	)

	if err != nil {
		return err
	}

	return nil
}

type loginAuth struct {
	username, password string
}

func buildAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("unknown from server")
		}
	}
	return nil, nil
}
