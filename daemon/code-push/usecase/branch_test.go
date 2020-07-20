package usecase_test

import (
	"errors"
	"github.com/funnyecho/code-push/daemon/code-push"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBranch(t *testing.T) {
	t.Run("invalid params", func(t *testing.T) {
		{
			branch, err := useCase.GetBranch(nil)
			assert.Error(t, err)
			assert.Nil(t, branch)
		}

		{
			branch, err := useCase.CreateBranch(nil)
			assert.Error(t, err)
			assert.Nil(t, branch)
		}

		{
			err := useCase.DeleteBranch(nil)
			assert.Error(t, err)
		}

		{
			token, err := useCase.GetBranchEncToken(nil)
			assert.Error(t, err)
			assert.Nil(t, token)
		}
	})

	branchName := []byte("code-push branch testing")
	var branchId []byte

	t.Run("fetch not existed branch", func(t *testing.T) {
		t.Log("if branchId not existed, no error occur, and return empty data")

		branch, err := useCase.GetBranch([]byte("not existed branch id"))
		assert.NoError(t, err)
		assert.Nil(t, branch)
	})

	t.Run("fetch encToken of not existed branch", func(t *testing.T) {
		t.Log("if branchId not existed, error occur")

		token, err := useCase.GetBranchEncToken([]byte("not existed branch id"))
		assert.True(t, errors.Is(err, code_push.ErrBranchNotFound))
		assert.Nil(t, token)
	})

	t.Run("create branch", func(t *testing.T) {
		branch, err := useCase.CreateBranch(branchName)
		assert.NoError(t, err)
		assert.NotNil(t, branch)
		assert.Equal(t, branch.Name, string(branchName))
		assert.True(t, len(branch.EncToken) > 0)

		branchId = []byte(branch.ID)

		t.Run("fetch existed branch", func(t *testing.T) {
			fetchBranch, fetchErr := useCase.GetBranch(branchId)
			assert.NoError(t, fetchErr)
			assert.NotNil(t, fetchBranch)
			assert.Equal(t, fetchBranch.Name, string(branchName))
		})

		t.Run("fetch encToken of existed branch", func(t *testing.T) {
			token, fetchErr := useCase.GetBranchEncToken(branchId)
			assert.NoError(t, fetchErr)
			assert.NotNil(t, token)
			assert.Equal(t, branch.EncToken, string(token))
		})
	})

	t.Run("create branch if branchName was existed", func(t *testing.T) {
		t.Log("Branch name is unique")

		branch, err := useCase.CreateBranch(branchName)
		assert.True(t, errors.Is(err, code_push.ErrBranchNameExisted))
		assert.Nil(t, branch)
	})

	t.Run("delete branch", func(t *testing.T) {
		deleteErr := useCase.DeleteBranch(branchId)
		assert.NoError(t, deleteErr)

		fetchBranch, fetchErr := useCase.GetBranch(branchId)
		assert.NoError(t, fetchErr)
		assert.Nil(t, fetchBranch)

		token, err := useCase.GetBranchEncToken(branchId)
		assert.True(t, errors.Is(err, code_push.ErrBranchNotFound))
		assert.Nil(t, token)
	})
}

func createRandomBranch() *code_push.Branch {
	branch, _ := useCase.CreateBranch([]byte(uuid.NewV4().String()))
	return branch
}