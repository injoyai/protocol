package http

import (
	"bytes"
	"fmt"
	"github.com/injoyai/base/g"
	"github.com/injoyai/conv"
	"io"
	"net/http"
	"net/http/httputil"
	"strings"
)

type Request http.Request

func (this *Request) Bytes() g.Bytes {
	bs, _ := httputil.DumpRequest((*http.Request)(this), true)
	return bs
}

func (this *Request) AddHeader(key string, val ...string) *Request {
	for _, v := range val {
		this.Header.Add(key, v)
	}
	return this
}

func (this *Request) AddHeaders(m map[string][]string) *Request {
	for k, v := range m {
		this.AddHeader(k, v...)
	}
	return this
}

func (this *Request) SetHeader(key string, val ...string) *Request {
	this.Header.Del(key)
	for _, v := range val {
		this.Header.Add(key, v)
	}
	return this
}

func (this *Request) SetHeaders(m map[string][]string) *Request {
	for k, v := range m {
		this.SetHeader(k, v...)
	}
	return this
}

func (this *Request) AddCookie(cookies ...*http.Cookie) *Request {
	for _, v := range cookies {
		(*http.Request)(this).AddCookie(v)
	}
	return this
}

func (this *Request) GetCookies() []*http.Cookie {
	return (*http.Request)(this).Cookies()
}

func (this *Request) GetCookie(key string) (*http.Cookie, error) {
	return (*http.Request)(this).Cookie(key)
}

func (this *Request) SetUserAgent(s string) *Request {
	return this.SetHeader("User-Agent", s)
}

// SetReferer 设置Referer
func (this *Request) SetReferer(s string) *Request {
	return this.SetHeader("Referer", s)
}

func (this *Request) SetAuthorization(s string) *Request {
	return this.SetHeader("Authorization", s)
}

func (this *Request) SetToken(s string) *Request {
	return this.SetHeader("Authorization", s)
}

func (this *Request) SetContentType(s string) *Request {
	return this.SetHeader("Content-Type", s)
}

func (this *Request) SetBody(body interface{}) *Request {
	switch val := body.(type) {
	case io.ReadCloser:
		this.Body = val
	case io.Reader:
		this.Body = io.NopCloser(val)
	default:
		this.Body = io.NopCloser(bytes.NewReader(conv.Bytes(body)))
	}
	return this
}

func NewRequest(method, url string, body io.Reader) *Request {
	r, _ := http.NewRequest(method, url, body)
	return (*Request)(r)
}

type Response struct {
	StatusCode     int
	StatusCodeText string
	Header         http.Header
	Body           []byte
}

func (this *Response) AddHeader(key, val string) *Response {
	this.Header[key] = append(this.Header[key], val)
	return this
}

func (this *Response) SetBody(body []byte) *Response {
	this.Body = body
	return this
}

func (this *Response) Bytes() g.Bytes {
	data := []byte(fmt.Sprintf("HTTP/1.1 %d %s\r\n", this.StatusCode, this.StatusCodeText))
	for i, v := range this.Header {
		data = append(data, []byte(fmt.Sprintf("%s: %s\r\n", i, strings.Join(v, "; ")))...)
	}
	data = append(data, []byte("\r\n")...)
	return append(data, this.Body...)
}

func NewResponse(statusCode int, body []byte) *Response {
	return &Response{
		StatusCode:     statusCode,
		StatusCodeText: "",
		Header: http.Header{
			"Content-Length": conv.Strings(len(body)),
			"Content-Type":   []string{"text/html", "charset=utf-8"},
		},
		Body: body,
	}
}

func NewResponseBytes(statusCode int, body []byte) g.Bytes {
	return NewResponse(statusCode, body).Bytes()
}

func NewResponseBytes200(body []byte) g.Bytes {
	return NewResponse(200, body).Bytes()
}

func NewResponseBytes204() g.Bytes {
	return NewResponse(204, nil).Bytes()
}

func NewResponseBytes301(addr string) g.Bytes {
	return NewResponse(301, nil).AddHeader("Location", addr).Bytes()
}

func NewResponseBytes302(addr string) g.Bytes {
	return NewResponse(302, nil).AddHeader("Location", addr).Bytes()
}

func NewResponse400(body []byte) g.Bytes {
	return NewResponseBytes(400, body)
}

func NewResponse500(body []byte) g.Bytes {
	return NewResponseBytes(500, body)
}
