package util

import "encoding/base64"

func EncodeBase64(value string) string {
	return base64.StdEncoding.EncodeToString([]byte(value))
}

func DecodeBase64(value string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(value)

	if err != nil {
		return "", err
	} else {
		return string(data), nil
	}
}
