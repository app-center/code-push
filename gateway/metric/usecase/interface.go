package usecase

type UseCase interface {
	RequestDuration
}

type RequestDuration interface {
	HttpDuration(svr, path string, success bool, durationSecond float64)
	GrpcDuration(svr, method string, success bool, durationSecond float64)
}
