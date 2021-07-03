package hy_http

import (
	"github.com/go-resty/resty/v2"
)

type Http interface {
	PostBytes(url string, headers map[string]string, data []byte) ([]byte, int, error)
	PostForm(url string, headers map[string]string, data map[string]string) ([]byte, int, error)
	PostJson(url string, headers map[string]string, data map[string]interface{}) ([]byte, int, error)
	Get(url string, headers map[string]string, data map[string]string) ([]byte, int, error)
}

type http struct {
}

func NewHttp() Http {
	return &http{}
}

func (h *http) PostBytes(url string, headers map[string]string, data []byte) ([]byte, int, error) {
	setReq := func(req *resty.Request) {
		req.SetHeader("Accept", "application/json")
		if len(data) > 0 {
			req.SetBody(data)
		}
	}
	return httpPostDo(url, headers, setReq)
}

func (h *http) PostForm(url string, headers map[string]string, data map[string]string) ([]byte, int, error) {
	setReq := func(req *resty.Request) {
		req.SetHeader("Content-Type", "application/x-www-form-urlencoded")
		if len(data) > 0 {
			req.SetFormData(data)
		}
	}
	return httpPostDo(url, headers, setReq)
}

func (h *http) PostJson(url string, headers map[string]string, data map[string]interface{}) ([]byte, int, error) {
	setReq := func(req *resty.Request) {
		req.SetHeader("Content-Type", "application/json")
		if len(data) > 0 {
			req.SetBody(data)
		}
	}
	return httpPostDo(url, headers, setReq)
}

func (h *http) Get(url string, headers map[string]string, data map[string]string) ([]byte, int, error) {
	setReq := func(req *resty.Request) {
		req.SetHeader("Content-Type", "application/json")
		if len(data) > 0 {
			req.SetQueryParams(data)
		}
	}
	return httpGetDo(url, headers, setReq)
}

func httpPostDo(url string, headers map[string]string, setReq func(req *resty.Request)) ([]byte, int, error) {
	client := resty.New()
	req := client.R()
	for k, v := range headers {
		req.SetHeader(k, v)
	}
	setReq(req)
	resp, err := req.Post(url)
	if err != nil {
		return nil, 0, err
	}
	return resp.Body(), resp.StatusCode(), nil
}

func httpGetDo(url string, headers map[string]string, setReq func(req *resty.Request)) ([]byte, int, error) {
	client := resty.New()
	req := client.R()
	for k, v := range headers {
		req.SetHeader(k, v)
	}
	setReq(req)
	resp, err := req.Get(url)
	if err != nil {
		return nil, 0, err
	}
	return resp.Body(), resp.StatusCode(), nil
}
