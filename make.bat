@echo off

@echo Enabling CGO...
@set CGO_ENABLED=1

@echo Building for Windows/386...
@set GOOS=windows
@set GOARCH=386
@go build -o clonehero-launcher-win-x86.exe

@echo Building for Windows/amd64...
@set GOOS=windows
@set GOARCH=amd64
@go build -o clonehero-launcher-win-x64.exe

@echo Done!
