package usecase

type UseCase interface {
	RequestDuration
}

type RequestDuration interface {
	GatewayDuration(svr, proto, path string, success bool, durationSecond float64)
	DaemonDuration(svr, proto, method string, success bool, durationSecond float64)
}
