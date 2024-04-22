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
example
```shell
go run main.go -env dev -port 5684
```
will be used ./config/dev.yml
