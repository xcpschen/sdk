package dingtalk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type (
	Client struct {
		accessToken string
		dtConf      IDtConf
		ac          *AToken
	}
	Respone struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}
)

const (
	dingtalkAPIHost string = "https://oapi.dingtalk.com"

	defaultReqTimeOut int = 60
	//默认5秒超时
)

const (
	ReqSucess int = 0
)

type DtReq interface {
	Uri(accessToken string) string
	Method() string
	Data() (val interface{}, err error)
	SetRespone(msg []byte) error
}

type IDtConf interface {
	GetCorpIDAndSecret() (string, string)
	GetAppIDAndSecret() (string, string)
	GetAgentID() string
}

func NewClientWithAccessToken(accessToken string) *Client {
	return &Client{accessToken: accessToken}
}

func NewClientWithIDtConf(c IDtConf) *Client {
	return &Client{dtConf: c}
}

//Exec 执行操作
//dtReq 必须实现DtReq 接口
func (dtp *Client) Exec(dtReq DtReq) error {
	var url string
	if dtp.accessToken == "" {
		ac, err := dtp.getAccessToken()
		if err != nil {
			return err
		}
		url = fmt.Sprintf("%s%s", dingtalkAPIHost, dtReq.Uri(ac))
	} else {
		url = fmt.Sprintf("%s%s", dingtalkAPIHost, dtReq.Uri(dtp.accessToken))
	}
	method := dtReq.Method()
	param, err := dtReq.Data()
	if err != nil {
		return err
	}

	req, err := NewHttpRequest(url, method, param)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: time.Duration(defaultReqTimeOut) * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("网络请求错误:%s", err)
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("结果解析错误:%s", err)
	}
	if err := dtReq.SetRespone(result); err != nil {
		return err
	}
	return nil
}

func NewHttpRequest(url, method string, data interface{}) (req *http.Request, err error) {
	if data == nil {
		req, err = http.NewRequest(method, url, nil)
	} else {
		jsonStr, _ := json.Marshal(data)
		req, err = http.NewRequest(method, url, bytes.NewBuffer(jsonStr))
	}
	return
}

type AToken struct {
	Value       string `json:"access_token"`
	ExpriesIn   int64  `json:"expires_in"`
	ExpriesDate int64
	Respone
}

func (dtp *Client) getAccessToken() (string, error) {
	now := time.Now().Unix()
	if dtp.ac != nil && dtp.ac.ExpriesDate > now {
		return dtp.ac.Value, nil
	}
	appid, appsecret := dtp.dtConf.GetAppIDAndSecret()
	url := fmt.Sprintf("%s/gettoken?appkey=%s&appsecret=%s", dingtalkAPIHost, appid, appsecret)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {

		return "", fmt.Errorf("结果解析错误:%s", err)
	}
	dtp.ac = &AToken{}
	if err := json.Unmarshal(result, dtp.ac); err != nil {
		return "", err
	}
	if dtp.ac.ErrCode != ReqSucess {
		return "", fmt.Errorf("%s", dtp.ac.ErrMsg)
	}
	dtp.ac.ExpriesDate = time.Now().Unix() + dtp.ac.ExpriesIn

	return dtp.ac.Value, nil
}
