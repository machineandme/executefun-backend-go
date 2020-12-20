package main

import (
	py3 "github.com/DataDog/go-python3"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestPythonRunning(t *testing.T) {
	pySetup()
	py3.PyRun_SimpleString("import sys\nsys.testA = True")
	if py3.PySys_GetObject("testA") != py3.Py_True {
		t.Error("Cannot set system values.")
	}
	pyStop()
}

func TestJWT(t *testing.T) {
	got, err := makeToken(map[string]string{"hello": "world"})
	if err != nil {
		t.Error(err)
	}
	expect := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VyRGF0YSI6eyJoZWxsbyI6IndvcmxkIn19.Zj7L47bBHxTkASnNwvKyuKKRmqABqN7cQt-V_Cs1Jfk_Gx_a-HzwKRAvf7WLNxbxl3uWR5rs6iXn5GLXLuf1tJ54wwDcxH8tqwveWsAmWf-8aXfqkqEthp8-xd5u8d6cJCeadxDDkBnrO5HAbZzQScVh5gZOxvTkBsTZObtFql4"
	if got != expect {
		t.Errorf("Got %v, expected %v.", got, expect)
	}
}

func TestServer(t *testing.T) {
	pySetup()
	pyFile("pythons/hello.py")
	http.HandleFunc("/full", makeHandler(HandlerConfig{
		includeHeaders:    true,
		includeRefHeaders: true,
		includeBody:       true,
	}))
	http.HandleFunc("/strict", makeHandler(HandlerConfig{
		includeHeaders:    false,
		includeRefHeaders: false,
		includeBody:       false,
	}))
	server := &http.Server{Addr: "127.0.0.1:8080", Handler: nil}
	end := make(chan bool, 1)
	go server.ListenAndServe()
	go func() {
		time.Sleep(100 * time.Millisecond)
		resp, err := http.Get("http://127.0.0.1:8080/strict")
		if err != nil {
			t.Errorf("Error %v", err)
			t.Fail()
		}
		body, err := ioutil.ReadAll(resp.Body)
		got := string(body)
		expect := "{\"message\":\"hello\",\"echo\":{\"body\":\"\",\"user_data\":{},\"headers\":{},\"query\":{}}}"
		if got != expect {
			t.Errorf("Got %v, expected %v.", got, expect)
		}
		time.Sleep(20 * time.Millisecond)
		sendBody := strings.NewReader("Ping")
		resp, err = http.Post("http://127.0.0.1:8080/full", "text/plain", sendBody)
		if err != nil {
			t.Errorf("Error %v", err)
			t.Fail()
		}
		body, err = ioutil.ReadAll(resp.Body)
		got = string(body)
		expect = "{\"message\":\"hello\",\"echo\":{\"user_data\":{},\"headers\":{\"User-Agent\":\"Go-http-client/1.1\",\"Accept-Encoding\":\"gzip\"},\"query\":{}}}"
		if got != expect {
			t.Errorf("Got %v, expected %v.", got, expect)
		}
		end <- true
	}()
	<-end
	server.Close()
	pyStop()
}