package qHttp

import "github.com/tidwall/gjson"

type ResponseDto struct {
	Body []byte
	Err  error
}

func (r *ResponseDto) Result() ([]byte, error) {
	return r.Body, r.Err
}

// 只返回data字段，errno!=0时返回BizError (调用方可以解析此错误里的错误码，进行特殊处理）
func (r *ResponseDto) ProcessResult() (*gjson.Result, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	resMap := gjson.ParseBytes(r.Body)
	errno := int(resMap.Get("code").Int())
	if errno != 0 {
		if resMap.Get("msg").Exists() {
			return nil, BizError{
				Errno:   errno,
				Message: resMap.Get("msg").String(),
			}
		}
		if resMap.Get("message").Exists() {
			return nil, BizError{
				Errno:   errno,
				Message: resMap.Get("message").String(),
			}
		}
		return nil, BizError{
			Errno:   errno,
			Message: "unknown reason",
		}
	}
	result := resMap.Get("result")
	if result.Exists() {
		return &result, nil
	}
	data := resMap.Get("data")
	return &data, nil
}
