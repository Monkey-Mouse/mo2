package mo2utils

import (
	"os"
	"testing"
)

func TestIsEnvRelease(t *testing.T) {
	tests := []struct {
		name        string
		wantRelease bool
	}{
		{"CHECK", os.Getenv("GIN_MODE") == "release"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRelease := IsEnvRelease(); gotRelease != tt.wantRelease {
				t.Errorf("IsEnvRelease() = %v, want %v", gotRelease, tt.wantRelease)
			}
		})
	}
}
func ExampleProcessAllFiles() {
	ProcessAllFiles("./", "/dist", func(parameter ...string) {
		println(parameter[0], "\n", parameter[2])
		//for _,v :=range parameter{
		//	println(v)
		//}
	})
	// Output:
}
