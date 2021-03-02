package emailservice

import (
	"fmt"
	"mo2/mo2utils/mo2errors"
	"net/http"
	"net/smtp"
	"os"
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
