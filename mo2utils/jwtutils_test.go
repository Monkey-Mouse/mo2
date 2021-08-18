package mo2utils

import (
	"reflect"
	"testing"

	"github.com/Monkey-Mouse/mo2/dto"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGenerateJwtCode(t *testing.T) {
	type args struct {
		info dto.LoginUserInfo
	}
	info := dto.LoginUserInfo{ID: primitive.NewObjectID(), Name: "xx", Email: "aaa", Roles: []string{"ss", "ll"}}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"test jwt", args{info}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateJwtCode(tt.args.info)
			if u, _ := ParseJwt(got); !reflect.DeepEqual(u, info) {
				t.Errorf("GenerateJwtCode() = %v, want %v", u, info)
			}
		})
	}
}
