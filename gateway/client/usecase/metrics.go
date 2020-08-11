package usecase

func (uc *useCase) RequestDuration(path string, success bool, durationSecond float64) {
	uc.metrics.RequestDuration(path, success, durationSecond)
}
