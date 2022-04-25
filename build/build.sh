#!/bin/sh
echo "Copy environment file"
yes | cp -rf build/fusionpbx-api-env /root/go/env/fusionpbx-api-env
echo "Build go application"
GOOS=linux GOARCH=amd64 go build -o fusionpbx-go-api main.go
echo "Restart service"
systemctl restart fusionpbx-go-api
systemctl status fusionpbx-go-api
