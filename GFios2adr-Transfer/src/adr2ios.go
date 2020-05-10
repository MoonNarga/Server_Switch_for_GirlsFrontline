package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/elazarl/goproxy"
	"io/ioutil"
	"net/http"
	"strings"
)

func Req4iOSServerList() goproxy.ReqConditionFunc {
	return func(req *http.Request, ctx *goproxy.ProxyCtx) bool {
		return req.URL.Host == adrServerListHost && req.URL.Path == adrServerListPath && req.Method == "POST"
	}
}

func ReplaceServerListParams4iOSServer(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return req, nil // TODO examine the effect of this error
	}
	req.Body.Close()

	if strings.Contains(string(body), "channel=cn_mica") {
		//log.Printf("channel == cn_appstore")
		// replace host
		req.URL.Host = iOSServerListHost
		req.Host = iOSServerListHost
		req.Header.Set("Host", iOSServerListHost)
		// replace params in the request body
		bodyStr := string(body)
		bodyStr = strings.Replace(bodyStr, "channel=cn_mica", "channel=cn_appstore", -1)
		bodyStr = strings.Replace(bodyStr, "platformChannelId=GWGW", "platformChannelId=ios", -1)
		bodyStr = strings.Replace(bodyStr, "device=adr", "device=ios", -1)
		newBody := ioutil.NopCloser(bufio.NewReader(strings.NewReader(bodyStr)))
		req.ContentLength = int64(len([]byte(bodyStr)))
		req.Header.Set("Content-Length", fmt.Sprintf("%d", req.ContentLength))
		req.Body = newBody
	} else {
		//log.Printf("channel != ios_appstore")
		req.Body = ioutil.NopCloser(bufio.NewReader(bytes.NewReader(body)))
	}

	return req, nil
}