package mo2utils

import (
	"os"
	"reflect"
	"testing"
)

func Test_initKey(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skipf("Skip for ci")
		return
	}
	_ = os.Remove("mo2.secret")
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "test first create", wantErr: false},
		{name: "test read", wantErr: false},
		{name: "test change", wantErr: false},
	}
	i := 0
	for _, tt := range tests {
		prev := make([]byte, 16)
		copy(prev, key)
		t.Run(tt.name, func(t *testing.T) {
			if err := initKey(); (err != nil) != tt.wantErr {
				t.Errorf("initKey() error = %v, wantErr %v", err, tt.wantErr)
				if i == 1 && !reflect.DeepEqual(key, prev) {
					t.Errorf("key should not change!")
					_ = os.Remove("mo2.secret")
				} else if i == 2 && reflect.DeepEqual(key, prev) {
					t.Errorf("key should change!")
				}
			}
		})
		i++
	}
}
