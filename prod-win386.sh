#!/bin/bash
name="logsrv386"
export CGO_ENABLED=0 
export GOOS=windows 
export GOARCH=386 
go build -tags=prod -o $name.exe main.go