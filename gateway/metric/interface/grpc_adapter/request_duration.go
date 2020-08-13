package metricAdapter

import (
	"context"
	"github.com/funnyecho/code-push/gateway/metric/interface/grpc/pb"
)

func (c *Client) HttpRequestDuration(svrName, path string, success bool, durationSecond float64) {
	c.requestDurationClient.Http(context.Background(), &pb.HttpRequestDurationRequest{
		Svr:            svrName,
		Path:           path,
		Success:        success,
		DurationSecond: durationSecond,
	})
}

func (c *Client) GrpcRequestDuration(svrName, method string, success bool, durationSecond float64) {
	c.requestDurationClient.Grpc(context.Background(), &pb.GrpcRequestDurationRequest{
		Svr:            svrName,
		Method:         method,
		Success:        success,
		DurationSecond: durationSecond,
	})
}
