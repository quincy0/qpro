package qHttp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	neturl "net/url"
	"strconv"
	"strings"
	"time"

	"github.com/quincy0/qpro/qLog"

	"github.com/google/uuid"

	"github.com/tidwall/gjson"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.uber.org/zap"
)

// Get 返回结果json化，并只取了result或data
func Get(ctx context.Context, url string, params map[string]any, opts ...ReqParamsOption) *ResponseDto {
	req := &requestParamsDto{
		Path:    url,
		method:  methodGet,
		Params:  params,
		timeout: defaultTimeout,
	}
	if opts != nil {
		for _, opt := range opts {
			opt(req)
		}
	}
	return do(ctx, req)
}

// Post 返回结果json化，并只取了result或data
func Post(ctx context.Context, url string, params map[string]any, opts ...ReqParamsOption) *ResponseDto {
	req := &requestParamsDto{
		Path:    url,
		method:  methodPost,
		Params:  params,
		timeout: defaultTimeout,
	}
	for _, opt := range opts {
		opt(req)
	}

	return do(ctx, req)
}

// GetNew 返回结果转json，没有其他处理
//
//	func GetNew(ctx context.Context, url string, params map[string]any, opts ...ReqParamsOption) (*gjson.Result, error) {
//		req := &requestParamsDto{
//			Path:    url,
//			method:  methodGet,
//			Params:  params,
//			timeout: defaultTimeout,
//		}
//		if opts != nil {
//			for _, opt := range opts {
//				opt(req)
//			}
//		}
//		bodyByt, err := do(ctx, req)
//		return _processResultNew(bodyByt, err)
//	}
//
//	func PostNew(ctx context.Context, url string, params map[string]any, opts ...ReqParamsOption) (*gjson.Result, error) {
//		req := &requestParamsDto{
//			Path:    url,
//			method:  methodPost,
//			Params:  params,
//			timeout: defaultTimeout,
//		}
//		for _, opt := range opts {
//			opt(req)
//		}
//
//		bodyByt, err := do(ctx, req)
//		return _processResultNew(bodyByt, err)
//	}
func Delete(ctx context.Context, url string, params map[string]any, opts ...ReqParamsOption) *ResponseDto {
	req := &requestParamsDto{
		Path:    url,
		method:  methodDelete,
		Params:  params,
		timeout: defaultTimeout,
	}
	for _, opt := range opts {
		opt(req)
	}

	return do(ctx, req)
}

func do(ctx context.Context, req *requestParamsDto) *ResponseDto {
	requestId := uuid.New().ID()
	startRequestTime := time.Now().UnixMilli()
	qLog.TraceInfo(
		ctx,
		"requestStart",
		zap.Uint32("requestId", requestId),
		zap.Any("req", req),
		zap.String("startTime", time.Now().Format("2006-01-02 15:04:05.9999")),
	)
	var err error
	var request *http.Request
	switch req.method {
	case methodPost:
		switch req.contentType {
		case ContentTypeSSML:

		case ContentTypeFormData:
			postValue := neturl.Values{}
			for k, v := range req.Params {
				postValue.Set(k, fmt.Sprintf("%v", v))
			}
			request, err = http.NewRequest("POST", req.Path, strings.NewReader(postValue.Encode()))
			request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			request.Header.Add("Content-Length", strconv.Itoa(len(postValue.Encode())))
		default:
			postValue := make(map[string]interface{})
			for k, v := range req.Params {
				postValue[k] = v
			}
			postBody, _ := json.Marshal(postValue)
			request, err = http.NewRequest("POST", req.Path, bytes.NewReader(postBody))
			request.Header.Add("Content-Type", "application/json;charset=utf-8")
			request.Header.Add("Content-Length", strconv.Itoa(len(postBody)))
		}

	case methodGet:
		getValues := neturl.Values{}
		for k, v := range req.Params {
			getValues.Set(k, fmt.Sprintf("%v", v))
		}
		if len(getValues) > 0 {
			req.Path += "?"
			req.Path += getValues.Encode()
		}
		request, err = http.NewRequest(methodGet, req.Path, nil)

	case methodDelete:
		getValues := neturl.Values{}
		for k, v := range req.Params {
			getValues.Set(k, fmt.Sprintf("%v", v))
		}
		if len(getValues) > 0 {
			req.Path += "?"
			req.Path += getValues.Encode()
		}
		request, err = http.NewRequest(methodDelete, req.Path, nil)
	}

	if len(req.header) > 0 {
		for k, v := range req.header {
			if len(v) > 0 {
				request.Header.Set(k, v[0])
			}
		}
	}
	// 设置超时时间
	ctx, cancel := context.WithCancel(ctx)
	time.AfterFunc(time.Duration(req.timeout)*time.Millisecond, func() {
		cancel()
	})
	request = request.WithContext(ctx)

	// 发起http请求
	aeClient := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	resp, err := aeClient.Do(request)

	if err != nil {
		qLog.TraceError(
			ctx,
			"requestEnd",
			zap.Uint32("requestId", requestId),
			zap.Int64("timeDuration", time.Now().UnixMilli()-startRequestTime),
			zap.String("endTime", time.Now().Format("2006-01-02 15:04:05.9999")),
			zap.Error(err),
		)
		return &ResponseDto{nil, err}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		qLog.TraceInfo(
			ctx,
			"requestEnd",
			zap.Uint32("requestId", requestId),
			zap.String("path", req.Path),
			zap.Int("code", resp.StatusCode),
			zap.Int64("timeDuration", time.Now().UnixMilli()-startRequestTime),
			zap.String("endTime", time.Now().Format("2006-01-02 15:04:05.9999")),
		)
		return &ResponseDto{nil, errors.New(resp.Status)}
	}

	var bodyByt []byte
	bodyByt, err = io.ReadAll(resp.Body)
	if err != nil {
		qLog.TraceError(
			ctx,
			"requestEnd",
			zap.Uint32("requestId", requestId),
			zap.String("path", req.Path),
			zap.Int("code", resp.StatusCode),
			zap.Error(err),
			zap.String("endTime", time.Now().Format("2006-01-02 15:04:05.9999")),
		)
		return &ResponseDto{nil, err}
	}

	qLog.TraceInfo(
		ctx,
		"requestEnd",
		zap.Uint32("requestId", requestId),
		zap.String("path", req.Path),
		zap.String("body", string(bodyByt)),
		zap.Int64("timeDuration", time.Now().UnixMilli()-startRequestTime),
		zap.String("endTime", time.Now().Format("2006-01-02 15:04:05.9999")),
	)
	return &ResponseDto{bodyByt, nil}
}

// 只返回data字段，errno!=0时返回BizError (调用方可以解析此错误里的错误码，进行特殊处理）
func _processResult(body []byte, requestErr error) (*gjson.Result, error) {
	if requestErr != nil {
		return nil, requestErr
	}
	resMap := gjson.ParseBytes(body)
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

// 只返回data字段，errno!=0时返回BizError (调用方可以解析此错误里的错误码，进行特殊处理）
func _processResultNew(body []byte, requestErr error) (*gjson.Result, error) {
	if requestErr != nil {
		return nil, requestErr
	}
	resMap := gjson.ParseBytes(body)
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
	return &resMap, nil
}
