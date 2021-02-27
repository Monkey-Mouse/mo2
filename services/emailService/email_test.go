package emailservice

import (
	"testing"
	"time"
)

func TestQueueEmail(t *testing.T) {
	SetFrequencyLimit(1)
	type args struct {
		msg        []byte
		receivers  []string
		remoteAddr string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "testStart", args: args{remoteAddr: "xxxx"}, wantErr: false},
		{name: "testBlok", args: args{remoteAddr: "xxxx"}, wantErr: true},
		{name: "testNoBlock", args: args{remoteAddr: "xxxxx"}, wantErr: false},
		{name: "testReleaseBlock", args: args{remoteAddr: "xxxx"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "testReleaseBlock" {
				time.Sleep(time.Second)
			}
			if err := QueueEmail(tt.args.msg, tt.args.receivers, tt.args.remoteAddr); (err != nil) != tt.wantErr {
				t.Errorf("QueueEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
