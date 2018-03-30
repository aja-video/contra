package utils

import (
	"contra/src/configuration"
	"net/smtp"
	"strconv"
)

// SendEmail to user
// - Consider using a library.
func SendEmail(c *configuration.Config, subject, message string) error {
	// Set up authentication information.
	auth := smtp.PlainAuth(
		"",
		c.SMTPUser,
		c.SMTPPass,
		c.SMTPHost,
	)
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
