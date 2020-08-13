package usecase

import (
	"fmt"
)

func (uc *useCase) HttpDuration(svr, path string, success bool, durationSecond float64) {
	uc.requestDurationMetric.With(
		"svr_name", svr,
		"interface", "http",
		"path", path,
		"success", fmt.Sprint(success),
	).Observe(durationSecond)
}

func (uc *useCase) GrpcDuration(svr, method string, success bool, durationSecond float64) {
	uc.requestDurationMetric.With(
		"svr_name", svr,
		"interface", "grpc",
		"path", method,
		"success", fmt.Sprint(success),
	).Observe(durationSecond)
}
