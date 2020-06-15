package code_push_pb

//go:generate protoc --gogofaster_out=plugins=grpc:. common.proto branch.proto env.proto version.proto
