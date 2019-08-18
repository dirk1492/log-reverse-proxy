package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	port := flag.Int("port", 8080, "Port number")
	target := flag.String("target-url", "", "Url of the target service")

	flag.Parse()

	remote, err := url.Parse(*target)
	if err != nil {
		panic(err)
	}

	proxy := NewProxy(remote)
	http.HandleFunc("/", handler(proxy))
	err = http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		panic(err)
	}
}

func NewProxy(target *url.URL) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host

		//path := req.Header.Get("X-Original-URI")
		//req.URL.Path = path
		log.Println(req.Method + " " + req.URL.String())
		txt := toJSON(req.Header)
		log.Println(txt)
	}

	modify := func(res *http.Response) error {

		log.Println(res.Status)
		txt := toJSON(res.Header)
		log.Println(txt)

		if res.StatusCode > 200 && res.StatusCode < 300 {
			res.StatusCode = 200
		} else {
			res.StatusCode = 401
		}

		return nil
	}

	return &httputil.ReverseProxy{
		Director:       director,
		ModifyResponse: modify,
	}
}

func handler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		rsp := NewLogWriter(w)
		p.ServeHTTP(rsp, r)
	}
}

func toJSON(val interface{}) string {
	data, err := json.Marshal(val)
	if err != nil {
		log.Println(err)
		return ""
	}

	return string(data)
}
