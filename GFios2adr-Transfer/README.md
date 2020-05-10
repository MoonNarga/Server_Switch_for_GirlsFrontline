# GirlsFrontline-Transfer
## Build
```bash
go build -o ./bin/GirlsFrontline-Transfer ./src/*
```

## Run
```bash
cd bin
./GirlsFrontline-Transfer
```

## Usage
```bash
cd bin
./GirlsFrontline-Transfer -h
Usage: GirlsFrontline-Transfer [options]
  -block
        enable blocking specified urls
  -conf string
        path to the configuration file for blacklist urls
  -dst string
        destination platform, default adr (default "adr")
  -h    print help message
  -port int
        port on which the proxy server listens to (default 8080)
  -src string
        source platform, default ios (default "ios")
  -v    enable verbose output
```
如果要让玩家只能进入到看板界面的话，需要在启动时添加 `-block` 选项。另外针对不同的需求，比如想要执行安卓到iOS的跨服的话，请添加 `-src adr -dst ios` 作为启动选项。程序默认是执行iOS到安卓的跨服。