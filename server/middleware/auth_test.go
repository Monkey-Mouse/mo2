package middleware

import (
	"reflect"
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

func Test_handlerMap_Group(t *testing.T) {
	type args struct {
		relativPath string
		roles       []string
	}
	tests := []struct {
		name string
		h    handlerMap
		args args
		want handlerMap
	}{
		{name: "test handler not change origin prefix", h: handlerMap{PrefixPath: ""}, args: args{"/xxx", nil}, want: handlerMap{PrefixPath: ""}},
		{name: "test handler not change origin role", h: handlerMap{PrefixPath: "", Roles: nil}, args: args{"/xxx", []string{"xxxx"}}, want: handlerMap{PrefixPath: "", Roles: nil}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.h.Group(tt.args.relativPath); !reflect.DeepEqual(tt.h, tt.want) {
				t.Errorf("handlerMap.Group() should not change origin map's value! changed: %v, origin %v", tt.h, tt.want)
			}
		})
	}
}
