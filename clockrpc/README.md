### 打包命令
```shell
set GOARCH=amd64&& set GOOS=linux&& go build -ldflags="-s -w" -o "clockone_Client" .\clockrpc\client\clockClient.go
```

```shell
set GOARCH=amd64&& set GOOS=linux&& go build -ldflags="-s -w" -o "clockone_Server" .\clockrpc\server\clockServer.go
```

### 压缩命令
```shell
upx.exe -9 --brute ./clockone_Client
upx.exe -9 --brute ./clockone_Server
```