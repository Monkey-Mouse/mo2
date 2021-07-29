package database

import "testing"

func TestCreateAccountIndex(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "test create account index", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateAccountIndex(); (err != nil) != tt.wantErr {
				t.Errorf("CreateAccountIndex() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
