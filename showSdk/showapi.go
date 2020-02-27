package showSdk

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
)

type NormalReq struct {
	url        string
	client     *fasthttp.Client
	bodyParams url.Values
	headParams url.Values
	file       map[string]string
	timeOut    time.Duration
}

// ShowAPIRequest 用于请求官网
func ShowAPIRequest(reqUrl string, appid int, sign string) *NormalReq {
	values := make(url.Values)
	values.Set("showapi_appid", strconv.Itoa(appid))
	values.Set("showapi_sign", sign)
	return &NormalReq{reqUrl, &fasthttp.Client{}, values, make(url.Values), make(map[string]string), 3 * time.Second}
}

// NormalRequest 通用请求
func NormalRequest(reqUrl string) *NormalReq {
	values := make(url.Values)
	return &NormalReq{reqUrl, &fasthttp.Client{}, values, make(url.Values), make(map[string]string), 3 * time.Second}
}

func (request *NormalReq) AddTextParam(key, value string) {
	request.bodyParams.Set(key, value)
}

func (request *NormalReq) AddFileParma(key, fileName string) {
	request.file[key] = fileName
}

func (request *NormalReq) AddHeadParma(key, value string) {
	request.headParams.Set(key, value)
}

func (request *NormalReq) SetTimeOut(timeOut time.Duration) {
	request.timeOut = timeOut
}

// Get get请求
func (request *NormalReq) Get() (string, error) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(strings.TrimSpace(request.url) + "?" + request.bodyParams.Encode())
	for k, v := range request.headParams {
		req.Header.Set(k, v[0])
	}
	req.Header.SetMethod("GET")
	resp := fasthttp.AcquireResponse()
	err := request.client.DoTimeout(req, resp, request.timeOut)
	if err != nil {
		return "", err
	}
	bodyBytes := resp.Body()
	return string(bodyBytes), nil
}

// Post 请求包括文件上传部分
func (request *NormalReq) Post() (string, error) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(strings.TrimSpace(request.url))
	for k, v := range request.bodyParams {
		req.PostArgs().Add(k, v[0])
	}
	for k, v := range request.headParams {
		req.Header.Set(k, v[0])
	}
	req.Header.SetMethod("POST")
	resp := fasthttp.AcquireResponse()
	err := request.client.DoTimeout(req, resp, request.timeOut)
	if err != nil {
		return "", err
	}
	bodyBytes := resp.Body()
	return string(bodyBytes), nil
}

// ParseJson 解析json
func ParseJson(req string) (map[string]interface{}, error) {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(req), &data); err != nil {
		return nil, errors.New("showapi return body is nil")
	}
	return data, nil
}

// Base64 图片文件传base64
func Base64(fileName string) string {
	fileBase64, _ := ioutil.ReadFile(fileName)
	return base64.StdEncoding.EncodeToString(fileBase64)
}
