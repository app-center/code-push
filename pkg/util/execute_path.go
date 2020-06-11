package util

import (
	"github.com/kardianos/osext"
)

func GetExecutePath() (string, error) {
	return osext.ExecutableFolder()
}
