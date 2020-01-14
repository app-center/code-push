package errors

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestError(t *testing.T) {
	nErr := fmt.Errorf("native errors")

	var decodeFunc DecodeFunc = func(params DecodeParams) string {
		v, _ := json.Marshal(params)
		return string(v)
	}

	err := New(CtorConfig{
		Error: nErr,
		Code:  "FA_ERROR",
		Msg:   "Unknown Error",
		Meta: MetaFields{
			"foo": 1,
		},
		DecodeFunc: decodeFunc,
	})
	expectErrOutput := decodeFunc(DecodeParams{
		Code: "FA_ERROR",
		Msg:  "Unknown Error",
		Meta: MetaFields{
			"foo": 1,
		},
	})

	if err.Error() != expectErrOutput {
		t.Fatalf(`Invalid errors output; expected: %s; got: %s`, expectErrOutput, err.Error())
	}

	if err.String() != expectErrOutput {
		t.Fatalf(`Invalid errors output; expected: %s; got: %s`, expectErrOutput, err.String())
	}

	if err.Unwrap() != nErr {
		t.Fatalf(`Invalid downstream output; expected: %v; got: %v`, nErr, err.Unwrap())
	}
}
