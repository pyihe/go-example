#!/usr/bin/env sh

protoc --go_out=. --go-grpc_out=. *.proto
