package main

import (
	"flag"
	"fmt"
	"github.com/elazarl/goproxy"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

const (
	//Listen            = 8888
	//Src               = "iOS"
	//Target            = "Android"
	blackListUrlsPath = "../conf/blacklist.json"
	//iOSServerListUrl  = "ios.transit.gf.ppgame.com/index.php"
	iOSServerListHost = "gfcn-transit.ios.sunborngame.com"
	iOSServerListPath = "/index.php"
	adrServerListPath = "/index.php"
	adrServerListHost = "gfcn-transit.gw.sunborngame.com"
	adrServerListOldHost = "adr.transit.gf.ppgame.com"
)

var (
	listenPort        int
	blackListConfFile string
	src               string
	dst               string
	block             bool
	verbose           bool
	needHelp          bool
)

func init() {
	flag.IntVar(&listenPort, "port", 8080, "port on which the proxy server listens to")
	flag.StringVar(&blackListConfFile, "conf", "", "path to the configuration file for blacklist urls")
	flag.StringVar(&src, "src", "ios", "source platform, default ios")
	flag.StringVar(&dst, "dst", "adr", "destination platform, default adr")
	flag.BoolVar(&verbose, "v", false, "enable verbose output")
	flag.BoolVar(&block, "block", false, "enable blocking specified urls")
	flag.BoolVar(&needHelp, "h", false, "print help message")
	flag.Parse()

	if verbose {
		fmt.Printf("listenPort: %d\nblackListConfFile:%s\nsrc:%s\ndst:%s\nverbose:%v\nban:%v\n",
			listenPort,
			blackListConfFile,
			src,
			dst,
			verbose,
			block)
	}
}

func main() {
	if needHelp {
		fmt.Println("Usage: GirlsFrontline-Transfer [options]")
		flag.PrintDefaults()
		os.Exit(0)
	}
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = verbose
	var blockedURLs *[]URL
	src = strings.ToLower(src)
	dst = strings.ToLower(dst)

	if block {
		if src == "ios" && dst == "adr" {
			if blackListConfFile == "" {
				err := GenerateBlackList(blackListUrlsPath, dst)
				if err != nil {
					log.Fatalf("Failed to generate configuration file for blacklist urls: %v", err)
				}
				blockedURLs = LoadConfiguration(blackListUrlsPath)
			} else {
				blockedURLs = LoadConfiguration(blackListConfFile)
			}

		} else if src == "adr" && dst == "ios" {
			if blackListConfFile == "" {
				err := GenerateBlackList(blackListUrlsPath, dst)
				if err != nil {
					log.Fatalf("Failed to generate configuration file for blacklist urls: %v", err)
				}
				blockedURLs = LoadConfiguration(blackListUrlsPath)
			} else {
				blockedURLs = LoadConfiguration(blackListConfFile)
			}
		} else {
			fmt.Printf("src 与 dst 不匹配，请指定正确的 src 以及 dst 参数。\n比如 src: ios, dst: adr。\n")
		}

		proxy.OnRequest(FilterRequests(blockedURLs)).DoFunc(
			func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
				//log.Println("Not allowed")
				return r, goproxy.NewResponse(r,
					goproxy.ContentTypeText,
					http.StatusForbidden,
					"Not allowed")
			})
	}

	if src == "ios" && dst == "adr" {
		proxy.OnRequest(Req4AdrServerList()).DoFunc(ReplaceServerListParams4AdrServer)
		//proxy.OnResponse(Resp4ServerList()).DoFunc(ReplaceRespOfServerList)
		//proxy.OnResponse(Resp4ClientVersion()).DoFunc(ReplaceClientVersion)
	} else if src == "adr" && dst == "ios"{
		proxy.OnRequest(Req4iOSServerList()).DoFunc(ReplaceServerListParams4iOSServer)
	} else {
		fmt.Printf("src 与 dst 不匹配，请指定正确的 src 以及 dst 参数。\n比如 src: ios, dst: adr。\n")
	}

	localIP := GetLocalIP()
	fmt.Printf("监听本机地址: %s, 端口: %d\n", localIP, listenPort)
	fmt.Printf("Girsfrontline-transfer %s -> %s 启动\n", src, dst)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", listenPort), proxy))
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