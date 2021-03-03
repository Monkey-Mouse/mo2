package controller

import (
	"io/ioutil"
	"strings"
	"testing"
)

func Test_GenUploadToken(t *testing.T) {
	req, resp := get(t, "/api"+imgGenToken, nil)
	r.ServeHTTP(resp, req)
	if resp.Code >= 300 {
		t.Errorf("bad response code")
		return
	}
	if p, err := ioutil.ReadAll(resp.Body); err != nil {
		t.Errorf("response err")
	} else {
		if strings.Contains(string(p), "Error") {
			t.Errorf("header response shouldn't return error: %s", p)
		} else if !strings.Contains(string(p), `token`) {
			t.Errorf("header response doen't match:\n%s", p)
		}
	}
}
