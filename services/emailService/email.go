package emailservice

import (
	"crypto/tls"
	"fmt"
	"mo2/mo2utils/mo2errors"
	"net/http"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/modern-go/concurrent"
	"github.com/willf/bloom"
)

// emailProp struct for send email
type emailProp struct {
	msg       []byte
	receivers []string
}

var emailChan chan<- emailProp
var initialed = false
var blockMap = concurrent.NewMap()
var bmChan = make(chan *concurrent.Map, 0)
var sec int64 = 5
var max int64 = 10
var blockTime int = 3600
var blockFilter = bloom.NewWithEstimates(10000, 0.01)

// SetFrequencyLimit set shortest resend time
func SetFrequencyLimit(seconds int64, limit int64, blocksec int) {
	sec = seconds
	max = limit
	blockTime = blocksec
}

func getBM() (bm *concurrent.Map) {
	select {
	case bm = <-bmChan:
	default:
		bm = blockMap
	}
	return
}

// QueueEmail add email to send queue
func QueueEmail(msg []byte, receivers []string, remoteAddr string) (err *mo2errors.Mo2Errors) {
	bm := getBM()
	if !initialed {
		startEmailService()
	}
	if blockFilter.TestString(remoteAddr) {
		err = mo2errors.New(http.StatusForbidden, "IP blocked! 检测到此IP潜在的ddos行为")
		return
	}
	val, ok := bm.Load(remoteAddr)
	prop := emailProp{msg: msg, receivers: receivers}
	if !ok {
		bm.Store(remoteAddr, int64(1))
		emailChan <- prop
		return
	}
	num := val.(int64)
	if num >= max {
		err = mo2errors.New(http.StatusTooManyRequests, "请求次数过多")
		blockFilter.AddString(remoteAddr)
		return
	}
	bm.Store(remoteAddr, num+1)
	emailChan <- prop
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
	go blockReseter()
	emailChan = emailc
	initialed = true
	return
}
func cleaner() {
	seconds := time.Second * time.Duration(sec)
	for {
		time.Sleep(seconds)
		nm := concurrent.NewMap()
		blockMap = nm
		bmChan <- nm
	}
}
func blockReseter() {
	seconds := time.Second * time.Duration(blockTime)
	for {
		time.Sleep(seconds)
		blockFilter.ClearAll()
	}
}
func startWorker(emailChan <-chan emailProp) {
	if os.Getenv("TEST") == "TRUE" {
		return
	}
	from := os.Getenv("emailAddr")
	password := os.Getenv("emailPass")
	// Sender data.

	// smtp server configuration.
	smtpHost := "smtpdm.aliyun.com"
	smtpPort := "465"
	addr := smtpHost + ":" + smtpPort
	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpHost,
	}
	fromH := []byte(fmt.Sprintf("From: %s\r\n", from))
	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)
	for {
		email := <-emailChan
		toH := []byte(fmt.Sprintf("To: %s\r\n", strings.Join(email.receivers, ",")))
		// Here is the key, you need to call tls.Dial instead of smtp.Dial
		// for smtp servers running on 465 that require an ssl connection
		// from the very beginning (no starttls)
		conn, err := tls.Dial("tcp", addr, tlsconfig)
		if err != nil {
			fmt.Println(err)
		}
		c, err := smtp.NewClient(conn, smtpHost)
		if err != nil {
			fmt.Println(err)
		}

		// Auth
		if err = c.Auth(auth); err != nil {
			fmt.Println(err)
		}

		// To && From
		if err = c.Mail(from); err != nil {
			fmt.Println(err)
		}

		for _, v := range email.receivers {
			if err = c.Rcpt(v); err != nil {
				fmt.Println(err)
			}
		}

		// Data
		w, err := c.Data()
		if err != nil {
			fmt.Println(err)
		}
		_, err = w.Write(fromH)
		if err != nil {
			fmt.Println(err)
		}
		_, err = w.Write(toH)
		if err != nil {
			fmt.Println(err)
		}
		_, err = w.Write(email.msg)
		if err != nil {
			fmt.Println(err)
		}

		err = w.Close()
		if err != nil {
			fmt.Println(err)
		}
		c.Quit()
	}
}
