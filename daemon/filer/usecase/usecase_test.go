package usecase_test

import (
	"github.com/funnyecho/code-push/daemon/filer"
	"github.com/funnyecho/code-push/daemon/filer/usecase"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"io"
	"os"
	"testing"
	"time"
)

var useCase *usecase.UseCase

func TestMain(m *testing.M) {
	adapters := &mockAdapters{
		files: make(map[string]*filer.File),
	}

	useCase = usecase.NewUseCase(usecase.CtorConfig{
		DomainAdapter: adapters,
		AliOssAdapter: adapters,
	})

	result := m.Run()

	os.Exit(result)
}

type mockAdapters struct {
	files map[string]*filer.File
}

func (m *mockAdapters) SignFetchURL(key []byte) ([]byte, error) {
	return key, nil
}

func (m *mockAdapters) Upload(stream io.Reader) ([]byte, error) {
	return []byte(uuid.NewV4().String()), nil
}

func (m *mockAdapters) File(fileKey string) (*filer.File, error) {
	return m.files[fileKey], nil
}

func (m *mockAdapters) InsertFile(file *filer.File) error {
	if file == nil {
		return errors.New("file required")
	}

	if file.Key == "" || file.Value == "" {
		return errors.New("file.Key and file.Value required")
	}

	fileToStorage := *file
	fileToStorage.CreateTime = time.Now()

	m.files[string(fileToStorage.Key)] = &fileToStorage
	return nil
}

func (m *mockAdapters) IsFileKeyExisted(fileKey string) bool {
	return m.files[fileKey] != nil
}
