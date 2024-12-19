## EVER-IOT DEVICE CLIENT APP

### Install

```shell
git clone https://github.com/ever-iot/device-go
cd ./device-go
go mod tidy
```

### Run

```shell
go run main.go -env <config-name> -port <port-number>
```
example for testing on linux device in dev mode will be used ./config/dev.yml
```
```shell
go run main.go -env dev
```

example for testing on Keenetic device in dev mode
```
```shell
go run main.go -env keenetic 
```

example for testing on Openwrt device in dev mode
```
```shell
go run main.go -env openwrt 
```
