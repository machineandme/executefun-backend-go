package main

import (
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestJWT(t *testing.T) {
	pySetup()
	pyFile("pythons/ts_user_data.py")
	h := makeHandler(HandlerConfig{
		includeHeaders:    false,
		includeRefHeaders: false,
		includeBody:       false,
	})

	req := httptest.NewRequest("GET", "https://execute.fun", nil)
	w := httptest.NewRecorder()
	h(w, req)
	resp := w.Result()
	body1, _ := ioutil.ReadAll(resp.Body)
	jwtToken1 := w.Header().Get("Set-Authorization")

	req = httptest.NewRequest("GET", "https://execute.fun?time=1", nil)
	req.Header.Set("Authorization", "Bearer "+jwtToken1)
	w = httptest.NewRecorder()
	h(w, req)
	resp = w.Result()
	body2, _ := ioutil.ReadAll(resp.Body)
	jwtToken2 := w.Header().Get("Set-Authorization")

	if jwtToken1 == jwtToken2 {
		t.Error("Tokens are same.")
	}
	if string(body1) == string(body2) {
		t.Error("Bodies are same.")
	}
	if !strings.Contains(string(body2), "\"user_data\":{\"time\":\"") {
		t.Error("User data not updated.")
	}

	pyStop()
}
