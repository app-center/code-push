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

	expectSign := util.EncodeMD5(fmt.Sprintf("%s.%s.%s", token, timestamp, nonce))

	return expectSign == sign, nil
}

func isTimestampValid(timestamp string) bool {
	startTime, startTimeErr := strconv.ParseInt(timestamp, 10, 64)

	if startTimeErr != nil {
		return false
	}

	nowTime := time.Now().UnixNano() / 1000000

	if startTime > nowTime {
		return false
	}

	return (nowTime - startTime) < authTimeout.Milliseconds()*1000*20
}
