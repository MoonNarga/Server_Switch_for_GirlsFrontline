package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/elazarl/goproxy"
	"io/ioutil"
	//"log"
	"net/http"
	"regexp"
	"strings"
)

func Req4AdrServerList() goproxy.ReqConditionFunc {
	return func(req *http.Request, ctx *goproxy.ProxyCtx) bool {
		//log.Printf("Method: %s, Host: %s, Path: %s\n", req.Method, req.URL.Host, req.URL.Path)
		return req.URL.Host == iOSServerListHost && req.URL.Path == iOSServerListPath && req.Method == "POST"
	}
}

func ReplaceServerListParams4AdrServer(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	//log.Printf("Start to replace server list params")

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return req, nil // TODO examine the effect of this error
	}
	req.Body.Close()

	if strings.Contains(string(body), "channel=cn_appstore") {
		//log.Printf("channel == cn_appstore")
		// replace host
		req.URL.Host = adrServerListHost
		req.Host = adrServerListHost
		req.Header.Set("Host", adrServerListHost)
		// replace params in the request body
		bodyStr := string(body)
		bodyStr = strings.Replace(bodyStr, "channel=cn_appstore", "channel=cn_mica", -1)
		bodyStr = strings.Replace(bodyStr, "platformChannelId=ios", "platformChannelId=GWGW", -1)
		bodyStr = strings.Replace(bodyStr, "device=ios", "device=adr", -1)
		//bodyStr = strings.Replace(bodyStr, "check_version=20500", "check_version=20501&device=adr", -1)
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

func Resp4ServerList() goproxy.RespConditionFunc {
	return func(resp *http.Response, ctx *goproxy.ProxyCtx) bool {
		//log.Printf("ctx.Req.URL.Host: %s, ctx.Req.Host: %s\n", ctx.Req.URL.Host, ctx.Req.Host)
		return ctx.Req.URL.Host == adrServerListHost && ctx.Req.URL.Path == adrServerListPath && ctx.Req.Method == "POST"
	}
}

func ReplaceRespOfServerList(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
	//log.Printf("Start to replace response of server list\n")
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp
	}
	resp.Body.Close()

	if strings.Contains(string(body), "<client_version>") {
		//log.Printf("<client_version> exists in the response body\n")
		newBody := string(body)
		reClientVersion := regexp.MustCompile(`<client_version>\d+</client_version>`)
		reTopClientVersion := regexp.MustCompile(`<top_client_version>\d+</top_client_version>`)

		newBody = reClientVersion.ReplaceAllString(newBody, "<client_version>20500</client_version>")
		newBody = reTopClientVersion.ReplaceAllString(newBody, "<top_client_version>2011</top_client_version>")

		resp.Body = ioutil.NopCloser(bufio.NewReader(strings.NewReader(newBody)))
		resp.ContentLength = int64(len([]byte(newBody)))
		resp.Header.Set("Content-Length", fmt.Sprintf("%d", resp.ContentLength))
	} else {
		resp.Body = ioutil.NopCloser(bufio.NewReader(bytes.NewReader(body)))
	}

	return resp
}

func Resp4ClientVersion() goproxy.RespConditionFunc {
	return func(resp *http.Response, ctx *goproxy.ProxyCtx) bool {
		return strings.Contains(ctx.Req.URL.Path, "/Index/version") && strings.HasPrefix(ctx.Req.URL.Host, "gfcn-game.gw.merge")
	}
}

func ReplaceClientVersion(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp
	}
	resp.Body.Close()

	if strings.Contains(string(body), "client_version") {
		re := regexp.MustCompile(`("client_version":)"\d+"`)
		bodyStr := string(body)
		bodyStr = re.ReplaceAllString(bodyStr, `$1"20500"`)

		resp.ContentLength = int64(len([]byte(bodyStr)))
		resp.Body = ioutil.NopCloser(bufio.NewReader(strings.NewReader(bodyStr)))
		resp.Header.Set("Content-Length", fmt.Sprintf("%d", resp.ContentLength))
	} else {
		resp.Body = ioutil.NopCloser(bufio.NewReader(bytes.NewReader(body)))
	}
	return resp
}

func CopyHeaderFromTo(src *http.Request, dst *http.Request) {
	for k, vs := range src.Header {
		for i, v := range vs {
			if i == 0 {
				dst.Header.Set(k, v)
			} else {
				dst.Header.Add(k, v)
			}
		}
	}
}

func FilterRequests(blockedURLs *[]URL) goproxy.ReqConditionFunc {
	return func(req *http.Request, ctx *goproxy.ProxyCtx) bool {
		requestUrl := req.URL.String()
		for _, blockedURL := range *blockedURLs {
			for _, blockedUrl := range blockedURL.GetUrls() {
				if strings.HasPrefix(requestUrl, blockedUrl) {
					return true // the requested url is in blacklist
				}
			}
		}
		return false
	}
}