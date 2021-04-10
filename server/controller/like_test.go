package controller

import (
	"encoding/json"
	"testing"

	"github.com/gin-gonic/gin"
)

func Benchmark_ginH(b *testing.B) {
	for n := 0; n < b.N; n++ {
		json.Marshal(gin.H{"test": true})
	}
}

var h = gin.H{"test": true}

func Benchmark_ginHPre(b *testing.B) {
	for n := 0; n < b.N; n++ {
		json.Marshal(h)
	}
}
