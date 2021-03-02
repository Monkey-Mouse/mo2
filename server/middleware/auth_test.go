package middleware

import (
	"fmt"
	"math/rand"
	"net/http"
	"path"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/modern-go/concurrent"
)

func Test_checkRL(t *testing.T) {
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

func Test_handlerMap_Get(t *testing.T) {
	type args struct {
		relativPath string
		handler     gin.HandlerFunc
		roles       []string
	}
	tests := []struct {
		name string
		h    handlerMap
		args args
	}{
		{name: "get test", h: H, args: args{"/xxx", nil, []string{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.h.Get(tt.args.relativPath, tt.args.handler, tt.args.roles...)
			_, ok := H.innerMap[handlerKey{tt.args.relativPath, http.MethodGet}]
			if !ok {
				t.Errorf("get test failed! failed to find handler after registered!")
			}
		})
	}
}

func Test_handlerMap_Handle(t *testing.T) {
	type args struct {
		method      string
		relativPath string
		handler     gin.HandlerFunc
		roles       []string
	}
	H = H.Group("/aaa")
	tests := []struct {
		name string
		h    handlerMap
		args args
	}{
		{name: "handler test", h: H, args: args{"xxx", "/xxx", nil, []string{}}},
		{name: "handler test", h: H, args: args{"xxaaa", "/ddxxx", nil, []string{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.h.Handle(tt.args.method, tt.args.relativPath, tt.args.handler, tt.args.roles...)
			_, ok := H.innerMap[handlerKey{path.Join("/aaa", tt.args.relativPath), tt.args.method}]
			if !ok {
				t.Errorf("handler test failed! failed to find handler after registered!")
			}
		})
	}
}

func Test_handlerMap_GetWithRL(t *testing.T) {
	h := handlerMap{handlers, "", make([][]string, 0), -1}
	type args struct {
		relativPath string
		handler     gin.HandlerFunc
		ratelimit   int
		roles       []string
	}
	tests := []struct {
		name string
		h    handlerMap
		args args
	}{
		{name: "get test", h: h, args: args{"/xxx", nil, 10, []string{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.h.GetWithRL(tt.args.relativPath, tt.args.handler, tt.args.ratelimit, tt.args.roles...)
			v, ok := h.innerMap[handlerKey{tt.args.relativPath, http.MethodGet}]
			if !ok {
				t.Errorf("get test failed! failed to find handler after registered!")
			}
			if v.limit != tt.args.ratelimit {
				t.Errorf("get test failed! rate limit value is wrong! expect: %v, real: %v", tt.args.ratelimit, v.limit)
			}
		})
	}
}

func Test_handlerMap_PostWithRL(t *testing.T) {
	h := handlerMap{handlers, "", make([][]string, 0), -1}
	type args struct {
		relativPath string
		handler     gin.HandlerFunc
		ratelimit   int
		roles       []string
	}
	tests := []struct {
		name string
		h    handlerMap
		args args
	}{
		{name: "post test", h: h, args: args{"/xxx", nil, 10, []string{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.h.PostWithRL(tt.args.relativPath, tt.args.handler, tt.args.ratelimit, tt.args.roles...)
			v, ok := h.innerMap[handlerKey{tt.args.relativPath, http.MethodGet}]
			if !ok {
				t.Errorf("post test failed! failed to find handler after registered!")
			}
			if v.limit != tt.args.ratelimit {
				t.Errorf("get test failed! rate limit value is wrong! expect: %v, real: %v", tt.args.ratelimit, v.limit)
			}
		})
	}
}
