package app

type Response[T any] struct {
	Code    int    `json:"code"`    // 错误码
	Data    T      `json:"data"`    // 返回数据
	Message string `json:"message"` // 错误信息
	Success bool   `json:"success"` // 是否成功
	TraceId string `json:"traceId"` // traceId
}

type Page[T any] struct {
	List      []T   `json:"list"`      //分页数据
	Count     int64 `json:"count"`     //总数
	PageIndex int   `json:"pageIndex"` //第几页
	PageSize  int   `json:"pageSize"`  //每页条数
}

func (res *Response[T]) ReturnOK() *Response[T] {
	res.Code = 0
	res.Success = true
	return res
}

func (res *Response[T]) ReturnError(code int, message string) *Response[T] {
	res.Code = code
	res.Message = message
	res.Success = false
	return res
}

type NullStruct struct {
}
