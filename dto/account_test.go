package dto

import "testing"

func TestContains(t *testing.T) {
	type args struct {
		slice []string
		item  string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "test contain", args: args{slice: []string{"xxxx", "xxx"}, item: "xxx"}, want: true},
		{name: "test not contain", args: args{slice: []string{"xxxx", "xxx"}, item: "xxxxx"}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Contains(tt.args.slice, tt.args.item); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}
