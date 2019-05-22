#!/usr/bin/env bash

protoc --go_out=plugins=grpc:. sys/pb/domainevents.proto

protoc --go_out=. chat/pb/events.proto

protoc --go_out=. iam/pb/events.proto

protoc --go_out=. photos/pb/events.proto