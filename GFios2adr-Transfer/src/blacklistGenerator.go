package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const (
	rankPath  = "/index.php/%d/Index/rank"
	homePath  = "/index.php/%d/Index/home"
	indexPath = "/index.php/%d/Index/index"
)

var (
	//blackListPaths = []string{rankPath, homePath, indexPath}
	blackListPaths = []string{rankPath, homePath}
)

//func init() {
//	flag.StringVar(&path, "path", "../conf/demo.json", "path to store the conf file")
//	flag.StringVar(&platform, "platform", "adr", "target platform for using the conf file")
//	flag.Parse()
//}

//func main() {
//	flag.StringVar(&path, "path", "../conf/demo.json", "path to store the conf file")
//	flag.StringVar(&platform, "platform", "adr", "target platform for using the conf file")
//	flag.Parse()
//
//	err := GenerateBlackList(path, platform)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Successfully generated blacklist of urls used for %s platform\n", platform)
//}

func GenerateBlackList(path string, targetPlatform string) error {
	var adrBlackListURLs []URL
	var iOSBlackListURLs []URL
	targetPlatform = strings.ToLower(targetPlatform)

	switch targetPlatform {
	case "ios":
		serverListIDs := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
		for _, serverListID := range serverListIDs {
			switch serverListID {
			case 0:
				iOSBlackListURLs = append(iOSBlackListURLs, URL{
					Scheme: "http",
					Host:   "gf-ios-cn-zs-game-0001.ppgame.com",
					Paths: []string{
						"/index.php/3000/Index/rank",
						"/index.php/3000/Index/home",
						//"/index.php/3000/Index/index",
					},
				})
			default:
				id := serverListID + 3000
				var paths []string
				for _, path := range blackListPaths {
					paths = append(paths, fmt.Sprintf(path, id))
				}
				iOSBlackListURLs = append(iOSBlackListURLs, URL{
					Scheme: "http",
					Host:   "s1.ios.gf.ppgame.com",
					Paths:  paths,
				})
			}
		}
	case "adr":
		serverListIDs := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}
		for _, serverListID := range serverListIDs {
			switch serverListID {
			case 0:
				adrBlackListURLs = append(adrBlackListURLs, URL{
					Scheme: "http",
					Host:   "gfcn-game.gw.merge.sunborngame.com",
					Paths: []string{
						"/index.php/1000/Index/rank",
						"/index.php/1000/Index/home",
						//"/index.php/1000/Index/index",
					},
				})
			default:
				id := serverListID + 1000
				var paths []string
				for _, path := range blackListPaths {
					paths = append(paths, fmt.Sprintf(path, id))
				}
				adrBlackListURLs = append(adrBlackListURLs, URL{
					Scheme: "http",
					Host:   "gfcn-game.gw.sunborngame.com",
					Paths:  paths,
				})
			}
		}
	default:
		return fmt.Errorf("wrong targetPlatform: %s", targetPlatform)
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	var data []byte
	if targetPlatform == "ios" {
		data, err = json.MarshalIndent(iOSBlackListURLs, "", "  ")
	} else {
		// adr
		data, err = json.MarshalIndent(adrBlackListURLs, "", "  ")
	}
	if err != nil {
		return err
	}
	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}
