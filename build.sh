export CGO_ENABLED=0
export GOOS=linux
#export GOOS=windows
#export GOOS=darwin
#export GOARCH=amd64
export GOARCH=arm64
go build http2mock.go vui.go apns.go