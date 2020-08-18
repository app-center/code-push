package usecase

func (uc *useCase) FileDownload(fileId []byte) ([]byte, error) {
	source, sourceErr := uc.filer.GetSource(fileId)
	if sourceErr != nil {
		return nil, sourceErr
	}

	return []byte(source.GetValue()), nil
}
