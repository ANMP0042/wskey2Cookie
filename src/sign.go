/**
 * @Author: YMBoom
 * @Description:
 * @File:  sign
 * @Version: 1.0.0
 * @Date: 2022/05/25 9:00
 */
package src

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"strings"
)

type Signer interface {
	sign() (data Result, err error)
}

func GetSign(singer Signer) (data Result, err error) {
	return singer.sign()
}

type Sign struct {
	body, functionId string
	Count            int
	conf             map[string]string
}

type RespData struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	Code      int    `json:"code"`
	Result    Result `json:"result"`
	Timestamp int64  `json:"timestamp"`
}

type Result struct {
	St            string `json:"st"`
	FunctionId    string `json:"functionId"`
	Sv            string `json:"sv"`
	Sign          string `json:"sign"`
	Client        string `json:"client"`
	Body          string `json:"body"`
	OriginalBody  string `json:"original_body"`
	ClientVersion string `json:"clientVersion"`
	Uuid          string `json:"uuid"`
}

func NewSign(functionId, body string) (iign *Sign, err error) {
	if functionId == "" || body == "" {
		err = errors.New("body or functionId is null")
		return
	}

	return &Sign{
		body:       body,
		functionId: functionId,
		conf:       GetMapStringMap("sign"),
	}, nil
}

func (s *Sign) sign() (data Result, err error) {
	values := url.Values{}
	values.Add("body", s.body)
	values.Add("functionId", s.functionId)
	payload := strings.NewReader(values.Encode())

	req := newReq(payload)

	u := fmt.Sprintf("%s%s%s", s.conf["gateway"], s.conf["port"], s.conf["path"])
	resp, err := req.SetUrl(u).
		SetMethod("POST").
		Do()

	if err != nil {
		return data, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return data, err
	}

	var res RespData
	if err = json.Unmarshal(body, &res); err != nil {
		return data, err
	}

	if !res.Success {
		err = errors.New(res.Message)
		return data, err
	}

	res.Result.OriginalBody = s.body
	return res.Result, nil
}

func newReq(payload io.Reader) *IRequest {
	headers := map[string]string{
		"cache-control": "no-cache",
		"content-type":  "application/x-www-form-urlencoded",
	}

	return NewIRequest(payload, headers)
}
