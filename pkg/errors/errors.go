package errors

import "fmt"

type MetaFields = map[string]interface{}
type EncodeParams struct {
	Code string
	Msg  string
	Meta MetaFields
}
type EncodeFunc func(params EncodeParams) string

func defaultErrorDecodeFunc(params EncodeParams) string {
	errMeta := ""

	if params.Meta != nil {
		errMeta = fmt.Sprintf("%+v", params.Meta)
	}

	return fmt.Sprintf("code-push errors: Code=%s\tMsg=%s\tMeta=%+v", params.Code, params.Msg, errMeta)
}

type metaError struct {
	err        error
	code       string
	msg        string
	meta       MetaFields
	encodeFunc EncodeFunc
}

func (err *metaError) Code() string {
	return err.code
}

func (err *metaError) Error() string {
	return err.encodeFunc(EncodeParams{
		Code: err.code,
		Msg:  err.msg,
		Meta: err.meta,
	})
}

func (err *metaError) String() string {
	return err.Error()
}

func (err *metaError) Unwrap() error {
	return err.err
}

type CtorConfig struct {
	Error      error
	Code       string
	Msg        string
	Meta       MetaFields
	EncodeFunc EncodeFunc
}

func Throw(config CtorConfig) *metaError {
	if config.EncodeFunc == nil {
		config.EncodeFunc = defaultErrorDecodeFunc
	}

	return &metaError{
		err:        config.Error,
		code:       config.Code,
		msg:        config.Msg,
		meta:       config.Meta,
		encodeFunc: config.EncodeFunc,
	}
}
