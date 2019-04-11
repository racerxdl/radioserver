#!/bin/bash

echo "Installing go-bindata"
go get -u github.com/go-bindata/go-bindata/...

echo "Building React App"
yarn build

echo "Bundling to ../webapp"
go-bindata -pkg webapp -o "../webapp/webapp.go" -prefix build/ build/*

echo "Done!"
