@echo off

@echo Building for Windows/386...
@set GOOS=windows
@set GOARCH=386
@go build -o clonehero-launcher-win-x86.exe

@echo Building for Windows/amd64...
@set GOOS=windows
@set GOARCH=amd64
@go build -o clonehero-launcher-win-x64.exe

@echo Building for Linux/386...
@set GOOS=linux
@set GOARCH=386
@go build -o clonehero-launcher-linux-x86_64

@echo Done!