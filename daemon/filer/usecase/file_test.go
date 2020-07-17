package usecase_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUseCase_GetSource(t *testing.T) {
	{
		file, err := useCase.GetSource(nil)
		assert.Error(t, err, "error when key is nil")
		assert.Nil(t, file, "nil file when key is nil")
	}

	{
		file, err := useCase.GetSource([]byte("not existed"))
		assert.Error(t, err)
		assert.Nil(t, file)
	}
}

func TestUseCase_InsertSource(t *testing.T) {
	key, err := useCase.InsertSource(nil, nil)
	assert.Error(t, err, "error when value is nil")
	assert.Nil(t, key, "nil key when error occur")

	key, err = useCase.InsertSource([]byte("no-scheme"), nil)
	assert.Error(t, err, "error when file value without scheme")
	assert.Nil(t, key, "nil key when error occur")

	key, err = useCase.InsertSource([]byte("ali-oss://val"), nil)
	assert.NoError(t, err, "`ali-oss` was supported")
	assert.NotNil(t, key, "valid key when scheme was supported")

	key, err = useCase.InsertSource([]byte("http://val"), nil)
	assert.NoError(t, err, "`http` was supported")
	assert.NotNil(t, key, "valid key when scheme was supported")

	key, err = useCase.InsertSource([]byte("https://val"), nil)
	assert.NoError(t, err, "`https` was supported")
	assert.NotNil(t, key, "valid key when scheme was supported")
}

func TestUseCase_Source_GetterSetter(t *testing.T) {
	fileDesc := "file desc"

	{
		fileVal := "ali-oss://val"
		fileKey, insertErr := useCase.InsertSource([]byte(fileVal), []byte(fileDesc))
		assert.NoError(t, insertErr)
		assert.NotNil(t, fileKey)

		fetchedFileVal, fetchErr := useCase.GetSource(fileKey)
		assert.NoError(t, fetchErr)
		assert.Equal(t, "val", string(fetchedFileVal), "ali-oss scheme getter will remove scheme")
	}

	{
		fileVal := "http://val"
		fileKey, insertErr := useCase.InsertSource([]byte(fileVal), []byte(fileDesc))
		assert.NoError(t, insertErr)
		assert.NotNil(t, fileKey)

		fetchedFileVal, fetchErr := useCase.GetSource(fileKey)
		assert.NoError(t, fetchErr)
		assert.Equal(t, fileVal, string(fetchedFileVal), "http scheme getter will return whole value")
	}

	{
		fileVal := "https://val"
		fileKey, insertErr := useCase.InsertSource([]byte(fileVal), []byte(fileDesc))
		assert.NoError(t, insertErr)
		assert.NotNil(t, fileKey)

		fetchedFileVal, fetchErr := useCase.GetSource(fileKey)
		assert.NoError(t, fetchErr)
		assert.Equal(t, fileVal, string(fetchedFileVal), "https scheme getter will return whole value")
	}
}
