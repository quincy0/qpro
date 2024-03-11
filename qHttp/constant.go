package qHttp

const (
	methodGet    = "GET"
	methodPost   = "POST"
	methodDelete = "DELETE"

	ContentTypeJson     = "json"
	ContentTypeFormData = "form-data"
	ContentTypeSSML     = "application/ssml+xml"

	defaultTimeout = 3000 // 单次请求超时时间
)
