#!/usr/bin/env bash

protoc --go_out=. chat/pb/protos.proto
protoc --go_out=plugins=grpc:. sys/rpc-domainevents/pb/protos.proto
protoc --go_out=plugins=grpc:. sys/scan/pb/protos.proto
protoc --go_out=. iam/pb/protos.proto
protoc --go_out=. photos/events/pb/photos.proto