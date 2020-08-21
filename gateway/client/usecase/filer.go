package usecase

import "context"

func (uc *useCase) FileDownload(ctx context.Context, fileId []byte) ([]byte, error) {
	source, sourceErr := uc.filer.GetSource(ctx, fileId)
	if sourceErr != nil {
		return nil, sourceErr
	}

	return []byte(source.GetValue()), nil
}
