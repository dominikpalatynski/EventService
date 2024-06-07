package main

import (
	"NotificationService/mail"
	"NotificationService/queue"
)

func main() {
	sender := mail.NewMailSender()
	queue.StartReceiving(sender)
}