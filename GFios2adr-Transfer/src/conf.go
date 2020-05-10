package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

const (
	//whiteList = "whitelist.json"
	//blackList = "blacklist.json"
)

var (
	//whiteListURLs []URL
	//blackListURLs []URL
)

type URL struct {
	Scheme string
	Host   string
	Paths  []string // each path in path must be a absolute path, namely, starts with "/"
}

func (u URL) GetUrls() []string {
	var res []string
	for _, path := range u.Paths {
		res = append(res, fmt.Sprintf("%s://%s%s", u.Scheme, u.Host, path))
	}

	return res
}

func LoadConfiguration(confFile string) *[]URL{
	cf, err := os.Open(confFile)
	if err != nil {
		log.Fatalf("loading conf: %v", err)
	}

	data, err := ioutil.ReadAll(cf)
	if err != nil {
		log.Fatalf("reading conf data: %v", err)
	}

	//if confFile == blackList {
	//	err = json.Unmarshal(data, &blackListURLs)
	//	if err != nil {
	//		log.Fatalf("unmarshalling json data: %v", err)
	//	}
	//} else if confFile == whiteList {
	//	err = json.Unmarshal(data, &whiteListURLs)
	//	if err != nil {
	//		log.Fatalf("unmarshalling json data: %v", err)
	//	}
	//}

	var res []URL
	err = json.Unmarshal(data, &res)
	if err != nil {
		log.Fatalf("unmarshalling json data: %v", err)
	}

	return &res
}
