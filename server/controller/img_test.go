package controller

import (
	"testing"
)

func Test_GenUploadToken(t *testing.T) {
	req := get(t, "/api"+apiImgGenToken, nil)
	req1 := req.Clone(req.Context())
	addCookie(req)
	ts := []tests{
		{name: "Test token gen", req: req, wantCode: 200, wantStr: "token"},
		{name: "Test auth", req: req1, wantCode: 403, wantStr: ""},
	}
	testHTTP(t, ts...)
}
