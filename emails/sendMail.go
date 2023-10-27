package emails

import (
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

// Send email to given recipients via SMTP. This function should be called as a go routine so that email processing happens in the background
func Send(recipients []string, subject string, body string) {
	message := gomail.NewMessage()

	fromEmail := os.Getenv("FROM_EMAIL")
	smtpServer := os.Getenv("SMTP_SERVER")
	smtpPortStr := os.Getenv("SMTP_PORT")
	smtpUsername := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	if fromEmail == "" || smtpServer == "" || smtpPortStr == "" || smtpUsername == "" || smtpPassword == "" {
		log.Errorln("Cannot send email. Email credentials are not connfigured in the env")
		return
	}

	message.SetHeader("From", fromEmail)
	message.SetHeader("To", recipients...)
	message.SetHeader("Subject", subject)
	message.SetBody("text/html", body)

	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		smtpPort = 587
	}

	dialer := gomail.NewDialer(smtpServer, smtpPort, smtpUsername, smtpPassword)

	if err := dialer.DialAndSend(message); err != nil {
		log.Errorln("Caught error while sending email: ", err)
		return
	}

	log.Infoln("Email sent successfully!")
}
