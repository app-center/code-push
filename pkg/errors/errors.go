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

type OpenError struct {
	err        error
	code       string
	msg        string
	meta       MetaFields
	encodeFunc EncodeFunc
}

type Error struct {
	*OpenError
}

func (err *OpenError) Error() string {
	return err.encodeFunc(EncodeParams{
		Code: err.code,
		Msg:  err.msg,
		Meta: err.meta,
	})
}

func (err *OpenError) String() string {
	return err.Error()
}

func (err *OpenError) Unwrap() error {
	return err.err
}

type CtorConfig struct {
	Error      error
	Code       string
	Msg        string
	Meta       MetaFields
	EncodeFunc EncodeFunc
}

func NewOpenError(config CtorConfig) *OpenError {
	if config.EncodeFunc == nil {
		config.EncodeFunc = defaultErrorDecodeFunc
	}

	return &OpenError{
		err:        config.Error,
		code:       config.Code,
		msg:        config.Msg,
		meta:       config.Meta,
		encodeFunc: config.EncodeFunc,
	}
}
