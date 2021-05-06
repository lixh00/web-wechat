package protocol

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"net/http"
)

// ReturnResponse Http请求的响应结构体封装
type ReturnResponse struct {
	*http.Response
	err error
}

// NewReturnResponse Constructor for ReturnResponse
func NewReturnResponse(response *http.Response, err error) *ReturnResponse {
	return &ReturnResponse{Response: response, err: err}
}

// Err 获取当前请求的错误
func (r *ReturnResponse) Err() error {
	return r.err
}

// ScanJSON json序列化
func (r *ReturnResponse) ScanJSON(v interface{}) error {
	if data, err := r.ReadAll(); err != nil {
		return err
	} else {
		return json.Unmarshal(data, v)
	}
}

// ScanXML xml序列化
func (r *ReturnResponse) ScanXML(v interface{}) error {
	if data, err := r.ReadAll(); err != nil {
		return err
	} else {
		return xml.Unmarshal(data, v)
	}
}

// ReadAll 读取请求体
func (r *ReturnResponse) ReadAll() ([]byte, error) {
	if r.Err() != nil {
		return nil, r.Err()
	}
	buffer := bytes.Buffer{}
	if _, err := buffer.ReadFrom(r.Body); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
