package usecase_test

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestUseCase_UploadToAliOss(t *testing.T) {
	fileKey, uploadErr := useCase.UploadToAliOss(nil)
	assert.Error(t, uploadErr)
	assert.Nil(t, fileKey)

	fileKey, uploadErr = useCase.UploadToAliOss(strings.NewReader("foo bar"))
	assert.NoError(t, uploadErr)
	assert.NotNil(t, fileKey)

	fileVal, fileValErr := useCase.GetSource(fileKey)
	assert.NoError(t, fileValErr)
	assert.NotNil(t, fileVal)
}
