package usecase_test

import (
	"errors"
	"github.com/funnyecho/code-push/daemon/code-push"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEnv(t *testing.T) {
	t.Run("invalid params", func(t *testing.T) {
		{
			env, err := useCase.GetEnv(nil)
			assert.Error(t, err)
			assert.Nil(t, env)
		}

		{
			env, err := useCase.CreateEnv(nil, nil)
			assert.Error(t, err)
			assert.Nil(t, env)

			env, err = useCase.CreateEnv([]byte("test"), nil)
			assert.Error(t, err)
			assert.Nil(t, env)

			env, err = useCase.CreateEnv(nil, []byte("test"))
			assert.Error(t, err)
			assert.Nil(t, env)
		}

		{
			err := useCase.DeleteEnv(nil)
			assert.Error(t, err)
		}

		{
			token, err := useCase.GetEnvEncToken(nil)
			assert.Error(t, err)
			assert.Nil(t, token)
		}
	})

	t.Run("create env on existed branch", func(t *testing.T) {
		branch, _ := useCase.CreateBranch([]byte("code-push env testing"))
		envName := []byte("android")

		env, createErr := useCase.CreateEnv([]byte(branch.ID), envName)
		assert.NoError(t, createErr)
		assert.NotNil(t, env)
		assert.Equal(t, branch.ID, env.BranchId)

		t.Run("create env if envName was existed", func(t *testing.T) {
			env, createErr := useCase.CreateEnv([]byte(branch.ID), envName)
			assert.True(t, errors.Is(createErr, code_push.ErrEnvNameExisted))
			assert.Nil(t, env)
		})

		t.Run("fetch existed env", func(t *testing.T) {
			fetchEnv, fetchErr := useCase.GetEnv([]byte(env.ID))
			assert.NoError(t, fetchErr)
			assert.NotNil(t, fetchEnv)
			assert.Equal(t, env.ID, fetchEnv.ID)
		})

		t.Run("fetch encToken of existed env", func(t *testing.T) {
			token, tokenErr := useCase.GetEnvEncToken([]byte(env.ID))
			assert.NoError(t, tokenErr)
			assert.NotNil(t, token)
			assert.Equal(t, env.EncToken, string(token))
		})

		t.Run("delete existed env", func(t *testing.T) {
			err := useCase.DeleteEnv([]byte(env.ID))
			assert.NoError(t, err)

			t.Run("create env if original env with same envName was deleted", func(t *testing.T) {
				env, createErr := useCase.CreateEnv([]byte(branch.ID), envName)
				assert.NoError(t, createErr)
				assert.NotNil(t, env)
				assert.Equal(t, branch.ID, env.BranchId)
			})

			t.Run("fetch deleted env", func(t *testing.T) {
				t.Log("if env was deleted, no error occur, and return empty data")
				fetchEnv, fetchErr := useCase.GetEnv([]byte(env.ID))
				assert.NoError(t, fetchErr)
				assert.Nil(t, fetchEnv)
			})

			t.Run("fetch encToken of deleted env", func(t *testing.T) {
				t.Log("if branchId not existed, error occur")
				token, tokenErr := useCase.GetEnvEncToken([]byte(env.ID))
				assert.True(t, errors.Is(tokenErr, code_push.ErrEnvNotFound))
				assert.Nil(t, token)
			})

			t.Run("delete env which was deleted", func(t *testing.T) {
				err := useCase.DeleteEnv([]byte(env.ID))
				assert.True(t, errors.Is(err, code_push.ErrEnvNotFound))
			})
		})
	})

	t.Run("create env on deleted branch", func(t *testing.T) {
		branch := createRandomBranch()
		envName := []byte("android")

		env, createErr := useCase.CreateEnv([]byte(branch.ID), envName)
		assert.NoError(t, createErr)
		assert.NotNil(t, env)
		assert.Equal(t, branch.ID, env.BranchId)

		_ = useCase.DeleteBranch([]byte(branch.ID))

		env, createErr = useCase.CreateEnv([]byte(branch.ID), []byte("ios"))
		assert.True(t, errors.Is(createErr, code_push.ErrBranchNotFound))
		assert.Nil(t, env)

	})

	t.Run("process on non existed env", func(t *testing.T) {
		envId := []byte("not existed env id")
		t.Run("fetch non existed env", func(t *testing.T) {
			t.Log("if env not existed, no error occur, and return empty data")
			fetchEnv, fetchErr := useCase.GetEnv(envId)
			assert.NoError(t, fetchErr)
			assert.Nil(t, fetchEnv)
		})

		t.Run("fetch encToken of deleted env", func(t *testing.T) {
			t.Log("if branchId not existed, error occur")
			token, tokenErr := useCase.GetEnvEncToken(envId)
			assert.True(t, errors.Is(tokenErr, code_push.ErrEnvNotFound))
			assert.Nil(t, token)
		})

		t.Run("delete env which was deleted", func(t *testing.T) {
			err := useCase.DeleteEnv(envId)
			assert.True(t, errors.Is(err, code_push.ErrEnvNotFound))
		})
	})
}

func createRandomEnv() (*code_push.Branch, *code_push.Env) {
	branch := createRandomBranch()
	env, _ := useCase.CreateEnv([]byte(branch.ID), []byte(uuid.NewV4().String()))
	return branch, env
}