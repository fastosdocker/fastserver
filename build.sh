go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct
go build main.go
./main