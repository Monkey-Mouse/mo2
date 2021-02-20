package mo2img

import (
	"regexp"
	"testing"
)

func TestGenerateUploadToken(t *testing.T) {
	tests := []struct {
		name      string
		wantToken string
	}{
		{name: "test token generation", wantToken: ".*"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotToken := GenerateUploadToken()

			if match, err := regexp.MatchString(tt.wantToken, gotToken); err != nil || !match {
				t.Errorf("GenerateUploadToken() = %v, want %v", gotToken, tt.wantToken)
			}
		})
	}
}
