package host

import (
	"fmt"
	"strconv"
	"strings"
	"wms-server/helpers"

	"github.com/go-gomail/gomail"
	logs "github.com/sirupsen/logrus"
)

type (
	mailSmpt struct {
		mailSmptDial *gomail.Dialer
		SMTP 		string
		Port		int
		EmailSender	string
		PassSender	string
		Log         *logs.Logger
	}
	// Host ...
	MailSmpt interface {
		SendEmailSMTP(to string, cc string, subject, message string) error 
	}
)

// SendEmailSMTP ..
func (ms *mailSmpt)SendEmailSMTP(to string, cc string, subject, message string) error {

	mailer := gomail.NewMessage()
	if cc != "" {
		mailer.SetHeaders(map[string][]string{
			"From": {ms.EmailSender},
			"To":      strings.Split(to,","),
			"Cc":     strings.Split(cc,","),
			"Subject": {subject},
		})
	} else {
		mailer.SetHeaders(map[string][]string{
			"From": {ms.EmailSender},
			"To":      strings.Split(to,","),
			"Subject": {subject},
		})
	}

	mailer.SetBody("text/html", message)

	err := ms.mailSmptDial.DialAndSend(mailer)
	if err != nil {
		return err
	}

	return nil
}

// InitializeMailGun ...
func InitializeMailSmpt(
	log *logs.Logger,
) MailSmpt {
	SMTP := helpers.GetEnv("SMTP_ADDRESS")
	Port := helpers.GetEnv("SMTP_PORT")
	EmailSender := helpers.GetEnv("EMAIL_SENDER")
	PassSender := helpers.GetEnv("EMAIL_PASSWORD")
	smtpPort,_ := strconv.Atoi(Port)
	fmt.Println(PassSender)
	dialer := gomail.NewDialer(SMTP, smtpPort, EmailSender, PassSender)
	return &mailSmpt{
		mailSmptDial: 	dialer,
		SMTP:			SMTP,
		Port:			smtpPort,
		EmailSender:	EmailSender,
		PassSender:		PassSender,
		Log:log,
	}
}
