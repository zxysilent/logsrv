#!/bin/bash
name="logsrv"
# export CGO_ENABLED=0 
# export GOOS=windows 
# export GOARCH=amd64 
go build -tags=prod -o $name.exe main.go


# Windows 下编译Linux 64位可执行程序

# SET GOOS=linux
# SET GOARCH=amd64

# go build
# GOOS：目标平台（darwin、freebsd、linux、windows） 

# GOARCH：目标平台的体系架构（386、amd64、arm）

# 交叉编译不支持 CGO

# window 后台方式运行

# go build -ldflags "-H=windowsgui"