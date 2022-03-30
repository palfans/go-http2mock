SET CGO_ENABLED=0
SET GOOS=linux
rem SET GOOS=windows
rem SET GOOS=darwin
rem SET GOARCH=arm64
SET GOARCH=amd64
go build http2mock.go vui.go apns.go