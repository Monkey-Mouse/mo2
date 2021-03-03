package controller

import (
	"mo2/dto"
	"mo2/mo2utils"
	"mo2/server/model"
	"testing"
)

func Test_GenUploadToken(t *testing.T) {
	req := get(t, "/api"+apiImgGenToken, nil)
	req1 := req.Clone(req.Context())
	req.Header.Set("Cookie",
		"jwtToken="+mo2utils.GenerateJwtCode(dto.LoginUserInfo{Roles: []string{model.OrdinaryUser}}))
	ts := []tests{
		{name: "Test token gen", req: req, wantCode: 200, wantStr: "token"},
		{name: "Test auth", req: req1, wantCode: 403, wantStr: ""},
	}
	testHTTP(t, ts...)
}
