package usecase_test

import (
	"errors"
	"github.com/funnyecho/code-push/daemon"
	"github.com/funnyecho/code-push/daemon/usecase"
	"github.com/funnyecho/code-push/pkg/log"
	uuid "github.com/satori/go.uuid"
	"io"
	"os"
	"testing"
	"time"
)

var useCase usecase.UseCase

func TestMain(m *testing.M) {
	adapters := &mockAdapters{
		branches: map[string]*daemon.Branch{},
		envs:     map[string]*daemon.Env{},
		versions: map[string]daemon.VersionList{},
	}

	useCase = usecase.NewUseCase(func(config *usecase.CtorConfig) {
		config.DomainAdapter = adapters
		config.AliOssAdapter = adapters
		config.Logger = log.NoopLogger
	})

	result := m.Run()

	os.Exit(result)
}

type mockAdapters struct {
	branches map[string]*daemon.Branch
	envs     map[string]*daemon.Env
	versions map[string]daemon.VersionList
	files 	 map[string]*daemon.File
}

func (m *mockAdapters) Branch(branchId []byte) (*daemon.Branch, error) {
	if branchId == nil {
		return nil, errors.New("branchId is required")
	}
	return m.branches[string(branchId)], nil
}

func (m *mockAdapters) CreateBranch(branch *daemon.Branch) error {
	if branch == nil {
		return errors.New("branch is required")
	}

	if _, existed := m.branches[branch.ID]; existed {
		return errors.New("branch was existed")
	}

	m.branches[branch.ID] = branch
	return nil
}

func (m *mockAdapters) DeleteBranch(branchId []byte) error {
	if branchId == nil {
		return errors.New("branchId is required")
	}

	if _, existed := m.branches[string(branchId)]; !existed {
		return errors.New("branch not found")
	}

	delete(m.branches, string(branchId))
	return nil
}

func (m *mockAdapters) IsBranchAvailable(branchId []byte) bool {
	if branchId == nil {
		return false
	}

	_, existed := m.branches[string(branchId)]
	return existed
}

func (m *mockAdapters) IsBranchNameExisted(branchName []byte) (bool, error) {
	if branchName == nil {
		return false, errors.New("branch name is required")
	}

	for _, b := range m.branches {
		if b.Name == string(branchName) {
			return true, nil
		}
	}

	return false, nil
}

func (m *mockAdapters) Env(envId []byte) (*daemon.Env, error) {
	if envId == nil {
		return nil, errors.New("envId is required")
	}

	return m.envs[string(envId)], nil
}

func (m *mockAdapters) CreateEnv(env *daemon.Env) error {
	if env == nil {
		return errors.New("env is required")
	}

	if _, existed := m.envs[env.ID]; existed {
		return errors.New("env was existed")
	}

	m.envs[env.ID] = env
	return nil
}

func (m *mockAdapters) DeleteEnv(envId []byte) error {
	if envId == nil {
		return errors.New("envId is required")
	}

	if _, existed := m.envs[string(envId)]; !existed {
		return errors.New("env not found")
	}

	delete(m.envs, string(envId))
	return nil
}

func (m *mockAdapters) IsEnvAvailable(envId []byte) bool {
	if envId == nil {
		return false
	}

	_, existed := m.envs[string(envId)]
	return existed
}

func (m *mockAdapters) IsEnvNameExisted(branchId, envName []byte) (bool, error) {
	if branchId == nil || envName == nil {
		return false, errors.New("branchId and envName are required")
	}

	for _, b := range m.envs {
		if b.BranchId == string(branchId) && b.Name == string(envName) {
			return true, nil
		}
	}

	return false, nil
}

func (m *mockAdapters) Version(envId, appVersion []byte) (*daemon.Version, error) {
	if envId == nil || appVersion == nil {
		return nil, errors.New("envId and appVersion are required")
	}

	versionList, existVersionList := m.versions[string(envId)]
	if versionList == nil || !existVersionList {
		return nil, nil
	}

	for _, ver := range versionList {
		if ver.AppVersion == string(appVersion) {
			return ver, nil
		}
	}

	return nil, nil
}

func (m *mockAdapters) VersionsWithEnvId(envId []byte) (daemon.VersionList, error) {
	if envId == nil {
		return nil, errors.New("envId is required")
	}

	return m.versions[string(envId)], nil
}

func (m *mockAdapters) CreateVersion(version *daemon.Version) error {
	if version == nil {
		return errors.New("version is required")
	}

	if available, _ := m.IsVersionAvailable([]byte(version.EnvId), []byte(version.AppVersion)); available {
		return errors.New("version was existed")
	}

	m.versions[version.EnvId] = append(m.versions[version.EnvId], version)
	return nil
}

func (m *mockAdapters) IsVersionAvailable(envId, appVersion []byte) (bool, error) {
	versionList, existVersionList := m.versions[string(envId)]
	if versionList == nil || !existVersionList {
		return false, nil
	}

	for _, ver := range versionList {
		if ver.AppVersion == string(appVersion) {
			return true, nil
		}
	}

	return false, nil
}

func (m *mockAdapters) SignFetchURL(key []byte) ([]byte, error) {
	return key, nil
}

func (m *mockAdapters) Upload(stream io.Reader) ([]byte, error) {
	return []byte(uuid.NewV4().String()), nil
}

func (m *mockAdapters) File(fileKey string) (*daemon.File, error) {
	return m.files[fileKey], nil
}

func (m *mockAdapters) InsertFile(file *daemon.File) error {
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
