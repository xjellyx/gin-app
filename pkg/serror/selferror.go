package serror

import (
	"bytes"

	"github.com/pkg/errors"
)

var _ SelfError = (*selfError)(nil)

type SelfError interface {
	// i 为了避免被其他包实现
	i()
	// Code 获取业务码
	Code() string
	Error() string
	StackError() error
}

type selfError struct {
	code     ErrorCode // 业务码
	message  string    // 错误描述
	stackErr error
}

func Error(businessCode ErrorCode, language string) SelfError {
	msg, err := LocalizeError(language, string(businessCode))
	if err != nil {
		panic(err)
	}
	return &selfError{
		code:    businessCode,
		message: msg,
	}
}

func (e *selfError) i() {}

func (e *selfError) Error() string {
	return e.message

}

// WithStack 堆栈
func (e *selfError) WithStack(err error) {
	e.stackErr = errors.WithStack(err)
}

// Code 错误代码
func (e *selfError) Code() string {
	return string(e.code)
}

func (e *selfError) Message() string {
	return e.message
}

// StackError 堆栈错误
func (e *selfError) StackError() error {
	return e.stackErr
}

type TranslateErr []string

func (t *TranslateErr) Error() string {
	var (
		str = bytes.NewBufferString("")
	)
	for i, v := range *t {
		_, _ = str.WriteString(v)
		if i != len(*t)-1 {
			_, _ = str.WriteString(", ")
		}
	}
	return str.String()
}
