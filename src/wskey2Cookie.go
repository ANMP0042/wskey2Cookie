// Package src /**
package src

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"
)

type Wskey2Cookie struct {
	wskey string
	conf  map[string]string
}

type tokenKey struct {
	Code     string `json:"code"`
	TokenKey string `json:"tokenkey"`
	Url      string `json:"url"`
}

//
func NewWskey2Cookie(wskey string) (w2c *Wskey2Cookie, err error) {
	if wskey == "" {
		err = errors.New("wskey or pin can not be null")
		writeErr("wskey or pin can not be null")
		return
	}

	return &Wskey2Cookie{
		wskey: wskey,
		conf:  GetMapStringMap("base"),
	}, nil
}

// 执行
func (w2c *Wskey2Cookie) Do() {
	sign, err := w2c.getSign()
	if err != nil {
		writeErr(err.Error())
		return
	}

	if sign.Sign == "" {
		err = errors.New("signature verification failed")
		writeErr("err: signature verification failed")
		return
	}

	tk, err := w2c.genToken(sign)
	if err != nil {
		return
	}

	if tk.Code != "0" {
		err = errors.New("tokenKey code is " + tk.Code)
		writeErr("tokenKey code is " + tk.Code)
		return
	}
	if tk.TokenKey == "" || tk.TokenKey == "xxx" {
		err = errors.New("tokenKey tokenKey is " + tk.TokenKey)
		writeErr("tokenKey tokenKey is " + tk.TokenKey)
		return
	}
	w2c.redirect(tk)
	return
}

// 跳转 获取响应头里面的pt_key
func (w2c *Wskey2Cookie) redirect(tk tokenKey) {
	reqUrl := fmt.Sprintf("%s?tokenKey=%s&to=%s", tk.Url, tk.TokenKey, w2c.conf["redirect_url"])

	resp, err := w2c.req(nil, nil).SetMethod("GET").SetUrl(reqUrl).Do()
	if err != nil {
		writeErr(err.Error())
		return
	}

	setCookie := resp.Header["Set-Cookie"]
	if setCookie == nil {
		writeErr("Set-Cookie not Found")
		return
	}

	ck := w2c.getCookie(setCookie)
	if ck == "" {
		writeErr("cookie is null")
		return
	}

	if strings.Index(ck, "app_open") == -1 {
		writeErr("wskey has expired")
		return
	}
	fmt.Println("cookie is written to cookie.log")
	writeSuc(ck)
	return
}

// 获取token
func (w2c *Wskey2Cookie) genToken(sign Result) (t tokenKey, err error) {
	values := url.Values{}
	values.Add("body", sign.OriginalBody)
	values.Add("client", sign.Client)
	values.Add("clientVersion", sign.ClientVersion)
	values.Add("sign", sign.Sign)
	values.Add("st", sign.St)
	values.Add("sv", sign.Sv)
	values.Add("uuid", sign.Uuid)

	headers := map[string]string{
		"Cookie":       w2c.wskey,
		"Content-Type": "application/x-www-form-urlencoded",
		"User-Agent":   "Dalvik/2.1.0 (Linux; U; Android 10; MIX 2S MIUI/V12.5.1.0.QDGCNXM)",
	}

	req := w2c.req(values, headers)

	reqUrl := fmt.Sprintf("%s?functionId=%s", w2c.conf["api_url"], w2c.conf["function_id"])
	resp, err := req.SetUrl(reqUrl).Do()
	if err != nil {
		writeErr(err.Error())
		return
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		writeErr(err.Error())
		return
	}

	var tk tokenKey
	if err = json.Unmarshal(b, &tk); err != nil {

	}
	if err != nil {
		return
	}
	return tk, nil
}

// 获取cookie
func (w2c *Wskey2Cookie) getCookie(setCookie []string) (ck string) {
	for _, cookie := range setCookie {
		values := strings.Split(cookie, ";")
		for _, value := range values {
			data := strings.Split(value, "=")
			if data[0] == "pt_key" || data[0] == "pt_pin" {
				ck += value + ";"
				continue
			}
		}
	}
	return
}

// 获取签名
func (w2c *Wskey2Cookie) getSign() (res Result, err error) {
	sign, err := NewSign(w2c.conf["function_id"], w2c.conf["body"])
	if err != nil {
		return
	}

	return GetSign(sign)
}

// 写入错误信息
func writeErr(msg string) {
	Write2File("err: "+msg, "cookie.txt")
}

// 写入cookie
func writeSuc(msg string) {
	Write2File("cookie: "+msg, "cookie.txt")
}

// 获取请求request
func (w2c *Wskey2Cookie) req(values url.Values, headers map[string]string) *IRequest {
	payload := strings.NewReader(values.Encode())
	return NewIRequest(payload, headers)
}
