// +build normal

package model

import (
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestAccount_IsValid(t *testing.T) {
	tests := []struct {
		name      string
		a         Account
		wantValid bool
	}{
		{name: "not valid test", a: Account{ID: primitive.NilObjectID}, wantValid: false},
		{name: "valid test", a: Account{ID: primitive.NewObjectID()}, wantValid: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotValid := tt.a.IsValid(); gotValid != tt.wantValid {
				t.Errorf("Account.IsValid() = %v, want %v", gotValid, tt.wantValid)
			}
		})
	}
}

func TestAddRoles(t *testing.T) {
	type args struct {
		a     *Account
		roles []Erole
	}
	acc := &Account{Roles: []Erole{"xxx", "yyyy"}}
	tests := []struct {
		name string
		args args
	}{
		{name: "test dup role", args: args{a: acc, roles: []Erole{"yyyy"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AddRoles(tt.args.a, tt.args.roles...)
			m := make(map[string]string, 0)
			for _, v := range acc.Roles {
				m[v] = v
			}
			if len(m) < len(acc.Roles) {
				t.Errorf("Unexpected dup role found!")
			}
		})
	}
}
