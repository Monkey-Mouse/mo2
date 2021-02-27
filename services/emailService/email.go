package emailservice

import (
	"errors"
	"fmt"
	"net/smtp"
	"os"
	"time"
)

// emailProp struct for send email
type emailProp struct {
	msg       []byte
	receivers []string
}

var emailChan chan<- emailProp
var initialed = false
var blockMap map[string]int64 = make(map[string]int64, 100)
var sec int64 = 30

// SetFrequencyLimit set shortest resend time
func SetFrequencyLimit(seconds int64) {
	sec = seconds
}

// QueueEmail add email to send queue
func QueueEmail(msg []byte, receivers []string, remoteAddr string) (err error) {
	if !initialed {
		startEmailService()
	}
	val, ok := blockMap[remoteAddr]
	prop := emailProp{msg: msg, receivers: receivers}
	if !ok {
		blockMap[remoteAddr] = time.Now().UnixNano()
		emailChan <- prop
		return
	}
	secs := (time.Now().Local().UnixNano() - val) / int64(time.Second)
	blockMap[remoteAddr] = time.Now().UnixNano()
	if secs >= sec {
		emailChan <- prop
		return
	}
	err = errors.New("Email请求过于频繁")
	return

}

// startEmailService start go routine for send email
func startEmailService() {
	if initialed {
		return
	}
	emailc := make(chan emailProp, 100)
	go startWorker(emailc)
	go cleaner()
	emailChan = emailc
	initialed = true
	return
}
func cleaner() {
	seconds := time.Second * time.Duration(sec)
	secondsint := int64(time.Second) * sec
	for {
		nano := time.Now().UnixNano()
		for k, v := range blockMap {
			if v < nano-secondsint {
				delete(blockMap, k)
			}
		}
		time.Sleep(seconds)
	}
}
func startWorker(emailChan <-chan emailProp) {
	from := os.Getenv("emailAddr")
	password := os.Getenv("emailPass")
	// Sender data.

	// smtp server configuration.
	smtpHost := "smtp.qq.com"
	smtpPort := "587"
	addr := smtpHost + ":" + smtpPort
	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)
	for {
		email := <-emailChan
		// Sending email.
		err := smtp.SendMail(addr, auth, from, email.receivers, email.msg)
		if err != nil {
			fmt.Println(err)
		}
	}
}
