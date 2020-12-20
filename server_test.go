package main

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	pySetup()
	pyFile("pythons/hello.py")
	http.HandleFunc("/ref", makeHandler(HandlerConfig{
		includeHeaders:    false,
		includeRefHeaders: true,
		includeBody:       false,
	}))
	http.HandleFunc("/strict", makeHandler(HandlerConfig{
		includeHeaders:    false,
		includeRefHeaders: false,
		includeBody:       false,
	}))
	http.HandleFunc("/body", makeHandler(HandlerConfig{
		includeHeaders:    false,
		includeRefHeaders: false,
		includeBody:       true,
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
		resp, err = http.Get("http://127.0.0.1:8080/ref")
		if err != nil {
			t.Errorf("Error %v", err)
			t.Fail()
		}
		body, err = ioutil.ReadAll(resp.Body)
		got = string(body)
		expect = "{\"message\":\"hello\",\"echo\":{\"body\":\"\",\"user_data\":{},\"headers\":{\"User-Agent\":\"Go-http-client/1.1\"},\"query\":{}}}"
		if got != expect {
			t.Errorf("Got %v, expected %v.", got, expect)
		}
		time.Sleep(20 * time.Millisecond)
		sendBody := strings.NewReader("Ping")
		resp, err = http.Post("http://127.0.0.1:8080/body", "text/plain", sendBody)
		if err != nil {
			t.Errorf("Error %v", err)
			t.Fail()
		}
		body, err = ioutil.ReadAll(resp.Body)
		got = string(body)
		expect = "{\"message\":\"hello\",\"echo\":{\"body\":\"Ping\",\"user_data\":{},\"headers\":{},\"query\":{}}}"
		if got != expect {
			t.Errorf("Got %v, expected %v.", got, expect)
		}
		end <- true
	}()
	<-end
	server.Close()
	pyStop()
}
