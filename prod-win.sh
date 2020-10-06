#!/bin/bash
swag init
name="logsrv"
go build -tags=prod -o $name.exe main.go


