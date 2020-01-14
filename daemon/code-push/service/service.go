package service

import (
	"github.com/funnyecho/code-push/daemon/code-push/repository"
	"log"
	"sync"
)

type CodePushService struct {
	mtx sync.RWMutex

	log *log.Logger

	branchRepo  *repository.IBranch
	envRepo     *repository.IEnv
	versionRepo *repository.IVersion
}
