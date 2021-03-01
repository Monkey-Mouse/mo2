package middleware

import (
	"fmt"
	"math/rand"
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
			if got := checkRL(tt.args.prop, tt.args.ip); got != tt.want {
				t.Errorf("checkRateLimit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handlerMap_Group1(t *testing.T) {
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
		{name: "test handler not change origin prefix", h: handlerMap{prefixPath: ""}, args: args{"/xxx", nil}, want: handlerMap{prefixPath: ""}},
		{name: "test handler not change origin role", h: handlerMap{prefixPath: "", roles: nil}, args: args{"/xxx", []string{"xxxx"}}, want: handlerMap{prefixPath: "", roles: nil}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.h.Group(tt.args.relativPath); !reflect.DeepEqual(tt.h, tt.want) {
				t.Errorf("handlerMap.Group() should not change origin map's value! changed: %v, origin %v", tt.h, tt.want)
			}
		})
	}
}

type testUser struct {
	roles []string
}

func (t testUser) IsInRole(role string) bool {
	for _, v := range t.roles {
		if v == role {
			return true
		}
	}
	return false
}
func Test_checkRoles(t *testing.T) {
	testU := testUser{roles: []string{"xx", "x"}}
	type args struct {
		uinfo        RoleHolder
		rolePolicies [][]string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "test or policy fail", args: args{uinfo: testU, rolePolicies: [][]string{{"xxx", "xxxx"}}}, wantErr: true},
		{name: "test or policy pass", args: args{uinfo: testU, rolePolicies: [][]string{{"xx", "xxxx"}}}, wantErr: false},
		{name: "test and policy fail", args: args{uinfo: testU, rolePolicies: [][]string{{"xx"}, {"xxxx"}}}, wantErr: true},
		{name: "test and policy pass", args: args{uinfo: testU, rolePolicies: [][]string{{"xx"}, {"x"}}}, wantErr: false},
		{name: "test combine policy pass", args: args{uinfo: testU, rolePolicies: [][]string{{"xxxxx", "xx"}, {"x", "pp"}}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkRoles(tt.args.uinfo, tt.args.rolePolicies); (err != nil) != tt.wantErr {
				t.Errorf("checkRoles() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_handlerMap_Group(t *testing.T) {
	type args struct {
		relativePath string
		roles        []string
	}
	tests := []struct {
		name string
		h    handlerMap
		args args
		want handlerMap
	}{
		{name: "test role policy and role", h: handlerMap{prefixPath: "/x", roles: [][]string{{"xx", "xxx"}}}, args: args{"xxx", []string{"xxxx"}}, want: handlerMap{prefixPath: "/x/xxx", roles: [][]string{{"xx", "xxx"}, {"xxxx"}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h.Group(tt.args.relativePath, tt.args.roles...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("handlerMap.Group() = %v, want %v", got, tt.want)
			}
		})
	}
}
func Benchmark_checkRoles(b *testing.B) {
	testU := testUser{roles: []string{"xxxxxxx", "xxxxxxxxxxxx"}}
	type args struct {
		uinfo        RoleHolder
		rolePolicies [][]string
	}
	rolePolicy := make([][]string, 100)
	ran := rand.Int63()
	for i := range rolePolicy {
		rolePolicy[i] = make([]string, 100)
		for u := range rolePolicy[i] {
			rolePolicy[i][u] = fmt.Sprintf("%v", ran)
			ran = (ran * 887) % 10000000
		}
	}
	arguments := args{
		uinfo:        testU,
		rolePolicies: rolePolicy,
	}
	b.ResetTimer()
	checkRoles(arguments.uinfo, arguments.rolePolicies)
}