go build -o tron_tools.exe  -trimpath -ldflags="-s -w" tron_tools.go

set GOOS=linux

go build -o tron_tools_linux -trimpath -ldflags="-s -w" tron_tools.go

set GOOS=darwin

go build -o tron_tools_darwin -trimpath -ldflags="-s -w" tron_tools.go