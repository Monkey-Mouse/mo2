package emailservice

import (
	"os"
	"testing"
	"time"
)

func TestQueueEmail(t *testing.T) {
	os.Setenv("TEST", "TRUE")
	SetFrequencyLimit(1, 3, 2)
	type args struct {
		msg        Mo2Email
		remoteAddr string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "testAdd1", args: args{remoteAddr: "xxxx"}, wantErr: false},
		{name: "testAdd2", args: args{remoteAddr: "xxxx"}, wantErr: false},
		{name: "testAdd3", args: args{remoteAddr: "xxxx"}, wantErr: false},
		{name: "testBlock", args: args{remoteAddr: "xxxx"}, wantErr: true},
		{name: "testNoBlock", args: args{remoteAddr: "xxxxx"}, wantErr: false},
		{name: "testReleaseBlock", args: args{remoteAddr: "xxxx"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "testReleaseBlock" {
				time.Sleep(time.Second * 3)
			}

			if err := QueueEmail(&tt.args.msg, tt.args.remoteAddr); (err != nil) != tt.wantErr {
				t.Errorf("QueueEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
