#!/usr/bin/env sh

Path=$pwd

#protoc --go_out=plugins=grpc:. $Path*.proto
protoc --go-grpc_out=plugins=grpc:. $Path*.proto
