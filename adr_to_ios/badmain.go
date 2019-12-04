package main

import (
	"bytes"
	// "crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"

	// "sync"
	"time"

	"github.com/ouqiang/goproxy"
)

// type Cache struct {
// 	m sync.Map
// }

// func (c *Cache) Set(host string, cert *tls.Certificate) {
// 	c.m.Store(host, cert)
// }
// func (c *Cache) Get(host string) *tls.Certificate {
// 	v, ok := c.m.Load(host)
// 	if !ok {
// 		return nil
// 	}

// 	return v.(*tls.Certificate)
// }

type EventHandler struct{}

func (e *EventHandler) Connect(ctx *goproxy.Context, rw http.ResponseWriter) {
	// 保存的数据可以在后面的回调方法中获取
	ctx.Data["req_id"] = "uuid"

	// 禁止访问某个域名
	if strings.Contains(ctx.Req.URL.Host, "example.com") {
		rw.WriteHeader(http.StatusForbidden)
		ctx.Abort()
		return
	}
}

func (e *EventHandler) Auth(ctx *goproxy.Context, rw http.ResponseWriter) {
	// 身份验证
}

func (e *EventHandler) BeforeRequest(ctx *goproxy.Context) {
	// 修改header
	ctx.Req.Header.Add("X-Request-Id", ctx.Data["req_id"].(string))
	// 设置X-Forwarded-For
	if clientIP, _, err := net.SplitHostPort(ctx.Req.RemoteAddr); err == nil {
		if prior, ok := ctx.Req.Header["X-Forwarded-For"]; ok {
			clientIP = strings.Join(prior, ", ") + ", " + clientIP
		}
		ctx.Req.Header.Set("X-Forwarded-For", clientIP)
	}
	if ctx.Req.URL.Hostname() == "adr.transit.gf.ppgame.com" {
		// ctx.Req.URL.Host = "ios.transit.gf.ppgame.com"
		// ctx.Req.Host = "ios.transit.gf.ppgame.com"
		ctx.Req.Header.Set("Host", "ios.transit.gf.ppgame.com")
		body, err := ioutil.ReadAll(ctx.Req.Body)
		if err != nil {
			// 错误处理
			return
		}
		// Request.Body只能读取一次, 读取后必须再放回去
		// Response.Body同理
		str := string(body)
		strings.Replace(str, "mica", "appstore", -1)
		body = []byte(str)
		ctx.Req.Body = ioutil.NopCloser(bytes.NewReader(body))
	}

	if ctx.Req.URL.Hostname() == "ios.transit.gf.ppgame.com" {
		// ctx.Req.URL.Host = "ios.transit.gf.ppgame.com"
		// ctx.Req.Host = "ios.transit.gf.ppgame.com"
		ctx.Req.Header.Set("Host", "adr.transit.gf.ppgame.com")
		body, err := ioutil.ReadAll(ctx.Req.Body)
		if err != nil {
			// 错误处理
			return
		}
		// Request.Body只能读取一次, 读取后必须再放回去
		// Response.Body同理
		str := string(body)
		strings.Replace(str, "appstore", "mica", -1)
		body = []byte(str)
		ctx.Req.Body = ioutil.NopCloser(bytes.NewReader(body))
	}
	// 读取Body

}

func (e *EventHandler) BeforeResponse(ctx *goproxy.Context, resp *http.Response, err error) {
	if err != nil {
		return
	}
	// 修改response
}

// 设置上级代理
func (e *EventHandler) ParentProxy(req *http.Request) (*url.URL, error) {
	// return url.Parse("http://127.0.0.1:8090")
	return nil, nil
	// return http.ProxyFromEnvironment(req)
}

func (e *EventHandler) Finish(ctx *goproxy.Context) {
	fmt.Printf("请求结束 URL:%s\n", ctx.Req.URL)
}

// 记录错误日志
func (e *EventHandler) ErrorLog(err error) {
	log.Println(err)
}

func main() {
	proxy := goproxy.New(goproxy.WithDelegate(&EventHandler{}))
	server := &http.Server{
		Addr:         ":8888",
		Handler:      proxy,
		ReadTimeout:  1 * time.Minute,
		WriteTimeout: 1 * time.Minute,
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
