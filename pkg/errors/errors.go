package errors

import "fmt"

type MetaFields = map[string]interface{}
type DecodeParams struct {
	Code string
	Msg  string
	Meta MetaFields
}
type DecodeFunc func(params DecodeParams) string

func defaultErrorDecodeFunc(params DecodeParams) string {
	errMeta := ""

	if params.Meta != nil {
		errMeta = fmt.Sprintf("%+v", params.Meta)
	}

	return fmt.Sprintf("code-push errors: Code=%s\tMsg=%s\tMeta=%+v", params.Code, params.Msg, errMeta)
}

type IError interface {
	Error() *Error
}

type Error struct {
	err        error
	code       string
	msg        string
	meta       MetaFields
	decodeFunc DecodeFunc
}

func (err *Error) Error() string {
	return err.decodeFunc(DecodeParams{
		Code: err.code,
		Msg:  err.msg,
		Meta: err.meta,
	})
}

func (err *Error) String() string {
	return err.Error()
}

func (err *Error) Unwrap() error {
	return err.err
}

type CtorConfig struct {
	Error      error
	Code       string
	Msg        string
	Meta       MetaFields
	DecodeFunc DecodeFunc
}

func New(config CtorConfig) *Error {
	if config.DecodeFunc == nil {
		config.DecodeFunc = defaultErrorDecodeFunc
	}

	return &Error{
		err:        config.Error,
		code:       config.Code,
		msg:        config.Msg,
		meta:       config.Meta,
		decodeFunc: config.DecodeFunc,
	}
}
