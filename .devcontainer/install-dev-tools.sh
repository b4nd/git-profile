#!/bin/bash

# Setup go 1.23.6
curl -sSfL https://dl.google.com/go/go1.23.6.linux-amd64.tar.gz -o go1.23.6.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.23.6.linux-amd64.tar.gz
rm go1.23.6.linux-amd64.tar.gz
sudo chown -R root:vscode /usr/local/go
sudo chmod -R 775 /usr/local/go

# Setup golangci-lint 
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sudo sh -s -- -b /usr/local/bin v1.42.1

# Setup gosec
curl -sfL https://raw.githubusercontent.com/securego/gosec/master/install.sh | sudo sh -s -- -b /usr/local/bin v2.22.0