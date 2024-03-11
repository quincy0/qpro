package qHttp

import (
	"errors"
	"fmt"
)

type BizError struct {
	Errno   int
	Message string
}

func (b BizError) Error() string {
	return fmt.Sprintf("errno[%d], errmsg:%s", b.Errno, b.Message)
}

// IsBizError 是否服务方错误码报错
func IsBizError(err error) (*BizError, bool) {
	bizErr := new(BizError)
	ok := errors.As(err, bizErr)
	if !ok {
		return nil, false
	}
	return bizErr, true
}
