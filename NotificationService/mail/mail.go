package mail

import (
	"NotificationService/util"
	"os"

	"gopkg.in/gomail.v2"
)

type MailSender struct {
	sender *gomail.Dialer
}

func newMessage(receiver string, title string) *gomail.Message {
	util.LoadEnv()

	message := gomail.NewMessage()

	message.SetHeader("From", os.Getenv("SENDER_NAME"))
	message.SetHeader("To", receiver)
	message.SetHeader("Subject", title)
	message.SetBody("text/plain", "test message body")

	return message
}

func NewMailSender() *MailSender {
	util.LoadEnv()
	sender := gomail.NewDialer(os.Getenv("SMTP_HOST"), 587, os.Getenv("SENDER_NAME"), os.Getenv("SENDER_PASSWORD"))

	return &MailSender{
		sender: sender,
	}
}

func (s *MailSender) SendMessage(title string) error {

	util.LoadEnv()

	message := newMessage(os.Getenv("SENDER_NAME"), title)

	return s.sender.DialAndSend(message)
}