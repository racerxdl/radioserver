#!/bin/sh

protoc -I ./ server.proto --go_out=plugins=grpc:.
