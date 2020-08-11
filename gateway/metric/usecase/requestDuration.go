package usecase

import (
	"fmt"
)

func (uc *useCase) GatewayDuration(svr, proto, path string, success bool, durationSecond float64) {
	uc.requestDurationMetric.With(
		"svr_type", "gateway",
		"svr_name", svr,
		"interface", proto,
		"path", path,
		"success", fmt.Sprint(success),
	).Observe(durationSecond)
}

func (uc *useCase) DaemonDuration(svr, proto, method string, success bool, durationSecond float64) {
	uc.requestDurationMetric.With(
		"svr_type", "gateway",
		"svr_name", svr,
		"interface", proto,
		"path", method,
		"success", fmt.Sprint(success),
	).Observe(durationSecond)
}
