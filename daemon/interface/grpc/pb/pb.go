package pb

//go:generate protoc --gogofaster_out=plugins=grpc:. common.proto
//go:generate protoc --gogofaster_out=plugins=grpc:. branch.proto env.proto version.proto
//go:generate protoc --gogofaster_out=plugins=grpc:. accessToken.proto
//go:generate protoc --gogofaster_out=plugins=grpc:. file.proto upload.proto
