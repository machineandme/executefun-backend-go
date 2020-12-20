package main

import (
	py3 "github.com/DataDog/go-python3"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)
func pyStop() {
	py3.Py_Finalize()
}
func pySetup() {
	py3.Py_Initialize()
	py3.PyRun_SimpleString("import sys")
}
func pyFile(filename string) {
	_, err := py3.PyRun_AnyFile(filename)
	if err != nil {
		panic(err)
	}
}

func setInterrupt() {
	lifeChannel := make(chan os.Signal, 1)
	signal.Notify(lifeChannel, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-lifeChannel
		log.Print("...")
		pyStop()
		os.Exit(0)
	}()
}

func main() {
	pySetup()
	setInterrupt()
	pyFile("pythons/hello.py")

	http.HandleFunc("/s", makeHandler(HandlerConfig{
		includeHeaders:    false,
		includeRefHeaders: false,
		includeBody:       false,
	}))
	http.HandleFunc("/f", makeHandler(HandlerConfig{
		includeHeaders:    true,
		includeRefHeaders: true,
		includeBody:       true,
	}))
	http.HandleFunc("/h", makeHandler(HandlerConfig{
		includeHeaders:    true,
		includeRefHeaders: false,
		includeBody:       false,
	}))
	http.HandleFunc("/r", makeHandler(HandlerConfig{
		includeHeaders:    false,
		includeRefHeaders: true,
		includeBody:       false,
	}))
	http.HandleFunc("/b", makeHandler(HandlerConfig{
		includeHeaders:    false,
		includeRefHeaders: false,
		includeBody:       true,
	}))

	log.Printf("runin: http://0.0.0.0:8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
