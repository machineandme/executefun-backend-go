package main

import (
	py3 "github.com/DataDog/go-python3"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var k_user_data = py3.PyUnicode_FromString("user_data")
var k_resp = py3.PyUnicode_FromString("response")

var BoringHeaders = []string{
	"Cache-Control",
	"Connection",
	"Expect",
	"Host",
	"Max-Forwards",
	"Pragma",
	"Range",
	"TE",
	"If-Match",
	"If-None-Match",
	"If-Modified-Since",
	"If-Unmodified-Since",
	"If-Range",
	"Accept",
	"Accept-Charset",
	"Accept-Encoding",
	"Accept-Language",

	"Content-Type",
	"Content-Length",

	"Upgrade-Insecure-Requests",
	"Vary",
	"Dnt",
}

var RefHeaders = []string{
	"From",
	"Referer",
	"User-Agent",
}

var AuthHeaders = []string{
	"Authorization",
	"Proxy-Authorization",
}

type HandlerConfig struct {
	includeHeaders    bool
	includeRefHeaders bool
	includeBody       bool
}

func processHeaders(conf HandlerConfig, request *http.Request, someDict *py3.PyObject, finished chan bool) {
	h := request.Header
	if !conf.includeHeaders {
		for _, i := range BoringHeaders {
			h.Del(i)
		}
	}
	if !conf.includeRefHeaders {
		for _, i := range RefHeaders {
			h.Del(i)
		}
	}
	var auth string
	auth = h.Get("Authorization")
	if auth == "" {
		auth = h.Get("Proxy-Authorization")
	}
	for _, i := range AuthHeaders {
		h.Del(i)
	}
	if auth != "" {
		authFields := strings.Fields(auth)
		if authFields[0] == "Bearer" {
			token := strings.Fields(auth)[1]
			userData := readToken(token)
			someDict.SetItem(
				py3.PyUnicode_FromString("user_data"),
				mapAsPyDict(userData),
			)
		} else {
			panic("Wrong formatted authorization header")
		}
	} else {
		someDict.SetItem(
			py3.PyUnicode_FromString("user_data"),
			py3.PyDict_New(),
		)
	}
	someDict.SetItem(
		py3.PyUnicode_FromString("headers"),
		maybeArrayMapAsPyDict(h),
	)
	finished <- true
}

func makeHandler(conf HandlerConfig) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		someDict := py3.PyDict_New()
		finished := make(chan bool)
		go processHeaders(conf, request, someDict, finished)
		someDict.SetItem(
			py3.PyUnicode_FromString("query"),
			maybeArrayMapAsPyDict(request.URL.Query()),
		)
		if conf.includeBody {
			bodyBytes, err := ioutil.ReadAll(request.Body)
			if err != nil {
				log.Fatal(err)
			}
			someDict.SetItem(
				py3.PyUnicode_FromString("body"),
				py3.PyUnicode_FromString(string(bodyBytes)),
			)
		}
		<-finished
		response := callSnake(someDict)
		tokenResponse, err := makeToken(pyMapToGo(response.GetItem(k_user_data)))
		fatal(err)
		writer.Header().Set("Set-Authorization", tokenResponse)
		resp := SerializePyObj(response.GetItem(k_resp))
		writer.Write([]byte(resp))
	}
}
