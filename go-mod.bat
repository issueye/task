rmdir /s/q vendor

set GOPROXY=https://goproxy.io,direct

set http_proxy=http://127.0.0.1:7897
set https_proxy=http://127.0.0.1:7897

:: 强制更新代码
go get -u github.com/issueye/ipc_grpc@v1.0.3

:: 更新依赖
go mod tidy
:: 更新本地依赖
go mod vendor