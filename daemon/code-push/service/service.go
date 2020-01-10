package service

import (
	"github.com/funnyecho/code-push/daemon/code-push/repository"
	"sync"
)

type CodePushService struct {
	mtx sync.RWMutex

	branchRepo *repository.IBranch
	envRepo *repository.IEnv
	versionRepo *repository.IVersion
}


