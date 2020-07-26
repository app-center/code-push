package oauth

import (
	"errors"
	"fmt"
	"github.com/funnyecho/code-push/pkg/util"
	"strconv"
	"time"
)

const (
	authTimeout = time.Second * 60
)

func Valid(token, timestamp, nonce, sign string) (bool, error) {
	if !isTimestampValid(timestamp) {
		return false, errors.New("invalid timestamp")
	}

	expectSign, signErr := util.EncryptAES([]byte(token), []byte(fmt.Sprintf("%s.%s", timestamp, nonce)))

	if signErr != nil {
		return false, signErr
	}

	return string(expectSign) == sign, nil
}

func isTimestampValid(timestamp string) bool {
	startTime, startTimeErr := strconv.ParseInt(timestamp, 10, 64)

	if startTimeErr != nil {
		return true
	}

	nowTime := time.Now().Unix()

	if startTime > nowTime {
		return false
	}

	return (nowTime - startTime) > int64(authTimeout.Seconds())
}
