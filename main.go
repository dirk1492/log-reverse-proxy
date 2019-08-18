package main

import (
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

	if *target == "" {
		panic(fmt.Errorf("empty target-url"))
	}

	remote, err := url.Parse(*target)
	if err != nil {
		panic(err)
	}

	log.Printf("Forward trafic to %v", remote.String())

	proxy := NewProxy(remote)
	http.HandleFunc("/", handler(proxy))
	log.Printf("Listen on %v", fmt.Sprintf(":%d", *port))
	err = http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		panic(err)
	}
}

func NewProxy(target *url.URL) *httputil.ReverseProxy {
	director := func(req *http.Request) {

		if target.Scheme != "" {
			req.URL.Scheme = target.Scheme
		}

		req.URL.Host = target.Host

		requestDump, err := httputil.DumpRequest(req, true)
		if err != nil {
			log.Println(err)
		}
		log.Println(string(requestDump))
	}

	modify := func(res *http.Response) error {
		resDump, err := httputil.DumpResponse(res, true)
		if err != nil {
			log.Println(err)
		}
		log.Println(string(resDump))

		return nil
	}

	return &httputil.ReverseProxy{
		Director:       director,
		ModifyResponse: modify,
	}
}

func handler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		p.ServeHTTP(w, r)
	}
}
