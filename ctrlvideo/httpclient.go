package ctrlvideo

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

// Hclient handle http request client
type Hclient struct {
	http.Client

	RetryTime     int
	RetryInterval time.Duration
}

// CommRequest 公共请求参数
type CommRequest struct {
	TimeStamp int64  `url:"timestamp"`
	AppID     string `url:"appid"`
	Nonce     string `url:"nonce"`
	Sign      string `url:"--"`
}

// NewCommRequest new CommRequest with appid string
func NewCommRequest(appID string) CommRequest {
	return CommRequest{
		TimeStamp: time.Now().Unix(),
		AppID:     appID,
		Nonce:     RandString(16),
	}
}

// ToHTTPReq get http request
func (c *CommRequest) ToHTTPReq(method string, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return req, nil
}

// CommRespone 公共响应参数
type CommRespone struct {
	Status int         `json:"status"`
	Msg    string      `json:"msg"`
	Msgzh  string      `json:"msgzh"`
	Data   interface{} `json:"data"`
}

// Parse Parse data
func (cr *CommRespone) Parse(obj interface{}) error {
	b, err := json.Marshal(cr.Data)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, obj)
}

// CommRespone statue data
const (
	ResponeFailureStatus int = iota
	ResponeSuccessStatus
)
const (
	// DefualtMaxRetry http request retry max time
	DefualtMaxRetry = 2
	//DefualtRetryInterval default value for retry interval time
	DefualtRetryInterval = 100 * time.Microsecond
)

// DoReq Do http request
func (h *Hclient) DoReq(r *http.Request) (respone *CommRespone, err error) {
	var resp *http.Response
	resp, err = h.doReq(r)
	if err != nil {
		return
	}
	var b []byte
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	resp.Body.Close()
	respone = new(CommRespone)
	if err = json.Unmarshal(b, respone); err != nil {
		err = fmt.Errorf("请求解析错误,返回信息,%s", string(b))
		return
	}
	if respone.Status != ResponeSuccessStatus {
		err = fmt.Errorf("操作失败，%s", respone.Msgzh+respone.Msg)
		return
	}
	return
}

// DoReqGetBigData when do request for big data
func (h *Hclient) DoReqGetBigData(r *http.Request, w io.WriteCloser) error {
	resp, err := h.doReq(r)
	if err != nil {
		return err
	}
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (h *Hclient) doReq(r *http.Request) (resp *http.Response, err error) {

	retry := 0
	if h.RetryTime == 0 {
		h.RetryTime = DefualtMaxRetry
	}
	if h.RetryInterval == 0 {
		h.RetryInterval = DefualtRetryInterval
	}
DoRetry:
	resp, err = h.Do(r)
	if err != nil {
		if retry >= h.RetryTime {
			err = fmt.Errorf("service is bad")
			return
		}
		retry++
		if h.RetryInterval > 0 {
			time.Sleep(h.RetryInterval)
		}
		goto DoRetry
	}
	return
}
