package database

import (
	"testing"
)

func TestIf(t *testing.T) {
	var a, b int = 2, 3
	res := If(a > b, a, b)
	if res == a {
		t.Error("三元表达式计算错误")
	}
	t.Log(res)
}
