echo Enabling CGO...
export CGO_ENABLED=1

echo Building for Linux/386...
export GOOS=linux
export GOARCH=386
go build -o clonehero-launcher-linux-x86

echo Building for Linux/amd64...
export GOOS=linux
export GOARCH=amd64
go build -o clonehero-launcher-linux-x64
