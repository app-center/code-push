package usecase_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUseCase_GetSource(t *testing.T) {
	{
		file, err := useCase.GetSource("")
		assert.Error(t, err, "error when key is nil")
		assert.Nil(t, file, "nil file when key is nil")
	}

	{
		file, err := useCase.GetSource("not existed")
		assert.Error(t, err)
		assert.Nil(t, file)
	}
}

func TestUseCase_InsertSource(t *testing.T) {
	key, err := useCase.InsertSource("", "", "", 0)
	assert.Error(t, err, "error when value is nil")
	assert.Nil(t, key, "nil key when error occur")

	key, err = useCase.InsertSource("no-scheme", "", "", 0)
	assert.Error(t, err, "error when file value without scheme")
	assert.Nil(t, key, "nil key when error occur")

	key, err = useCase.InsertSource("ali-oss://val", "", "", 0)
	assert.NoError(t, err, "`ali-oss` was supported")
	assert.NotNil(t, key, "valid key when scheme was supported")

	key, err = useCase.InsertSource("http://val", "", "", 0)
	assert.NoError(t, err, "`http` was supported")
	assert.NotNil(t, key, "valid key when scheme was supported")

	key, err = useCase.InsertSource("https://val", "", "", 0)
	assert.NoError(t, err, "`https` was supported")
	assert.NotNil(t, key, "valid key when scheme was supported")
}

func TestUseCase_Source_GetterSetter(t *testing.T) {
	fileDesc := "file desc"

	{
		fileVal := "ali-oss://val"
		fileKey, insertErr := useCase.InsertSource(fileVal, fileDesc, "", 0)
		assert.NoError(t, insertErr)
		assert.NotNil(t, fileKey)

		fetchedFileVal, fetchErr := useCase.GetSource(string(fileKey))
		assert.NoError(t, fetchErr)
		assert.Equal(t, "val", fetchedFileVal, "ali-oss scheme getter will remove scheme")
	}

	{
		fileVal := "http://val"
		fileKey, insertErr := useCase.InsertSource(fileVal, fileDesc, "", 0)
		assert.NoError(t, insertErr)
		assert.NotNil(t, fileKey)

		fetchedFileVal, fetchErr := useCase.GetSource(string(fileKey))
		assert.NoError(t, fetchErr)
		assert.Equal(t, fileVal, fetchedFileVal, "http scheme getter will return whole value")
	}

	{
		fileVal := "https://val"
		fileKey, insertErr := useCase.InsertSource(fileVal, fileDesc, "", 0)
		assert.NoError(t, insertErr)
		assert.NotNil(t, fileKey)

		fetchedFileVal, fetchErr := useCase.GetSource(string(fileKey))
		assert.NoError(t, fetchErr)
		assert.Equal(t, fileVal, fetchedFileVal, "https scheme getter will return whole value")
	}
}
