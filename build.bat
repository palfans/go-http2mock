SET CGO_ENABLED=0
SET GOOS=linux
rem SET GOOS=windows
SET GOARCH=amd64
go build http2mock.go vui.go apns.go