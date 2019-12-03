package main

import (
	"bytes"
	"fmt"
	"github.com/elazarl/goproxy"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

func main() {
	proxy := goproxy.NewProxyHttpServer()
	proxy.OnRequest(match()).DoFunc(serverList())
	log.Println("ProxyServer starts successfully")
	log.Printf("Listening on %s:%d\n", GetLocalIP(), 8888)
	log.Fatal(http.ListenAndServe(":8888", proxy))
}

func match() goproxy.ReqConditionFunc {
	return func(req *http.Request, ctx *goproxy.ProxyCtx) bool {
		if req.URL.Host == "adr.transit.gf.ppgame.com" {
			if strings.HasPrefix(req.URL.Path, "/index") || strings.HasPrefix(req.URL.Path, "index"){
				return true
			} else {
				return false
			}
		} else {
			return false
		}
	}
}

func serverList() func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	return func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		var url string
		if strings.HasPrefix(req.URL.Path, "/") {
			url = fmt.Sprintf("%s://%s%s", req.URL.Scheme, strings.Replace(req.URL.Host, "adr", "ios", -1), req.URL.Path)
		} else {
			url = fmt.Sprintf("%s://%s/%s", req.URL.Scheme, strings.Replace(req.URL.Host, "adr", "ios", -1), req.URL.Path)
		}
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Fatal(err)
			return nil, nil
		}
		data := strings.Replace(string(body), "mica", "appstore", -1)
		r, err := http.NewRequest(req.Method, url, bytes.NewBuffer([]byte(data)))
		if err != nil {
			log.Fatal(err)
			return nil, nil
		}
		setHeader(req, r)
		log.Println("Patch done.")
		return r, nil
	}
}

func setHeader(or, nr *http.Request) {
	for k, v := range or.Header {
		for _, vl := range v {
			nr.Header.Add(k, vl)
		}
	}
}

func GetLocalIP() string {
	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)
	found := false
	var rst string
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil && ipv4.String() != "127.0.0.1" {
			found = true
			rst = ipv4.String()
		}
	}
	if found {
		return rst
	}
	log.Fatal("No IPv4 addr found")
	os.Exit(1)
	return ""
}
