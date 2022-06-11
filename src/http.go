/**
 * @Author: YMBoom
 * @Description:
 * @File:  http
 * @Version: 1.0.0
 * @Date: 2022/05/25 9:00
 */
package src

import (
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type IRequest struct {
	Req    *http.Request
	Client *http.Client
	mu     sync.RWMutex
}

func NewIRequest(body io.Reader, headers map[string]string) *IRequest {
	req, _ := http.NewRequest("POST", "", body)
	iReq := &IRequest{
		Req: req,
		Client: &http.Client{
			Timeout: 60 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
			//CheckRedirect: myCheckRedirect,
		},
		mu: sync.RWMutex{},
	}

	if headers != nil {
		for i, v := range headers {
			iReq.AddHeader(i, v)
		}
	}
	return iReq
}

func (i *IRequest) AddHeader(name, value string) *IRequest {
	i.Req.Header.Add(name, value)
	return i
}

func (i *IRequest) SetHeader(name, value string) *IRequest {
	i.Req.Header.Set(name, value)
	return i
}

func (i *IRequest) SetCookie(name, value string) *IRequest {
	cookie := &http.Cookie{
		Name:  name,
		Value: value,
	}
	i.Req.AddCookie(cookie)
	return i
}

func (i *IRequest) SetCookies(cookies []*http.Cookie) *IRequest {
	for _, cookie := range cookies {
		i.Req.AddCookie(cookie)
	}
	return i
}
func (i *IRequest) SetMethod(method string) *IRequest {
	i.Req.Method = method
	return i
}

func (i *IRequest) SetUrl(reqUrl string) *IRequest {
	u, _ := url.Parse(reqUrl)
	i.Req.URL = u
	return i
}

func (i *IRequest) SetPostForm(val url.Values) *IRequest {
	i.Req.PostForm = val
	return i
}

func (i *IRequest) Do() (*http.Response, error) {
	return i.Client.Do(i.Req)
}
