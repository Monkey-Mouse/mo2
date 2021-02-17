package demo

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendEmail() {

	// Sender data.
	from := os.Getenv("emailAddr")
	password := os.Getenv("emailPass")
	// Receiver email address.
	to := []string{
		"liziyi@hust.edu.cn",
	}

	// smtp server configuration.
	smtpHost := "smtp.qq.com"
	smtpPort := "587"

	// Message.
	message := []byte("This is a test email message.")

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)
	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent Successfully!")
}
