package middleware

import (
	"testing"

	"github.com/modern-go/concurrent"
)

func Test_checkRateLimit(t *testing.T) {
	type args struct {
		prop handlerProp
		ip   string
	}
	prop := handlerProp{rates: concurrent.NewMap(), limit: 3}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "test add 1", args: args{prop: prop, ip: "xxxx"}, want: true},
		{name: "test add 2", args: args{prop: prop, ip: "xxxx"}, want: true},
		{name: "test add 3", args: args{prop: prop, ip: "xxxx"}, want: true},
		{name: "test block ip", args: args{prop: prop, ip: "xxxx"}, want: false},
		{name: "test another ip", args: args{prop: prop, ip: "xxxxy"}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkRateLimit(tt.args.prop, tt.args.ip); got != tt.want {
				t.Errorf("checkRateLimit() = %v, want %v", got, tt.want)
			}
		})
	}
}
