package errors

import (
	"fmt"
	assert "github.com/stretchr/testify/require"
	"testing"
)

func TestError(t *testing.T) {
	nErr := fmt.Errorf("native errors")

	encodeFunc := func(params EncodeParams) string {
		return fmt.Sprintf("Code: %s; Msg: %s; Meta.foo: %d", params.Code, params.Msg, params.Meta["foo"])
	}

	openParams := EncodeParams{
		Code: "FA_ERROR",
		Msg:  "Unknown Error",
		Meta: MetaFields{
			"foo": 1,
		},
	}

	expectErrMessage := fmt.Sprintf("Code: %s; Msg: %s; Meta.foo: %d", openParams.Code, openParams.Msg, openParams.Meta["foo"])

	err := NewOpenError(CtorConfig{
		Error:      nErr,
		Code:       openParams.Code,
		Msg:        openParams.Msg,
		Meta:       openParams.Meta,
		EncodeFunc: encodeFunc,
	})

	assert.Equal(t, expectErrMessage, err.Error())
	assert.Equal(t, expectErrMessage, err.String())
	assert.Equal(t, nErr, err.Unwrap())
}
