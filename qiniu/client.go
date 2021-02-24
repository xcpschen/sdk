package qiniu

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"time"
)

// Qclient http clientor
type Qclient struct {
	http.Client
	RetryTime     int
	RetryInterval time.Duration
}

const (
	// DefualtMaxRetry http request retry max time
	DefualtMaxRetry = 2
	//DefualtRetryInterval default value for retry interval time
	DefualtRetryInterval = 100 * time.Microsecond
)

// IReq http requestor interface
type IReq interface {
	// ToHTTPReq return a http.Request pointer and error when is failure
	ToHTTPReq() (req *http.Request, err error)
	// SaveRespone save respone form qiniu
	SaveRespone(QiniuRespone)
}

// IReqParam qiniu request param to string interace
type IReqParam interface {
	ToString() string
}

// DoReq Do http request
func (h *Qclient) DoReq(r IReq) (err error) {
	req, err := r.ToHTTPReq()
	if err != nil {
		return
	}
	retry := 0
	if h.RetryTime == 0 {
		h.RetryTime = DefualtMaxRetry
	}
	if h.RetryInterval == 0 {
		h.RetryInterval = DefualtRetryInterval
	}
DoRetry:
	resp, err := h.Do(req)
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
	var b []byte
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	resp.Body.Close()

	r.SaveRespone(NewQiniuRespone(b, resp.StatusCode, resp.Status))
	return
}

// QiniuRespone parse respone form http request
type QiniuRespone struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
	Body  []byte
}

func NewQiniuRespone(b []byte, statusCode int, status string) QiniuRespone {
	rq := QiniuRespone{}
	rq.Code = statusCode
	rq.Error = status
	json.Unmarshal(b, &rq)
	rq.Body = b
	return rq
}

// ParseBody 解析
func (rq QiniuRespone) ParseBody(obj interface{}) error {
	return json.Unmarshal(rq.Body, obj)
}

// SimpleReq 构建请求结构
type SimpleReq struct {
	// ReqType 请求类型
	ReqType     int
	Host        string
	URI         string
	Method      string
	ContentType string
	ak          string
	sk          string

	Param   IReqParam
	XQinius *XQiniu
	Resp    QiniuRespone
}

// SimpleReq 请求类型
const (
	ReqTypeForManger = iota
	ReqTypeForUpload
	ReqTypeForDownload
)

// ToHTTPReq implement IReq
func (s *SimpleReq) ToHTTPReq() (req *http.Request, err error) {
	req, err = http.NewRequest(s.Method, "https://"+s.Host+s.URI, s.setBody())
	if err != nil {
		return nil, err
	}
	if s.ContentType == "" {
		s.ContentType = "application/x-www-form-urlencoded"
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		req.Header.Set("Content-Type", s.ContentType)
	}
	if s.XQinius != nil {
		s.XQinius.SetToReqHeaders(req)
	}
	token := s.getToken()
	req.Header.Set("Authorization", token)
	return req, nil
}

// SetAkSK setting ak ,sk from implementing for IQiniuConf
func (s *SimpleReq) SetAkSK(c IQiniuConf) {
	s.ak, s.sk = c.GetAKSK()
}

// SetAkSK2 setting ak ,sk
func (s *SimpleReq) SetAkSK2(ak, sk string) {
	s.ak, s.sk = ak, sk
}

// SaveRespone save respone form qiniu
func (s *SimpleReq) SaveRespone(resp QiniuRespone) {
	s.Resp = resp
}

// setToken generate token for qiniu request
func (s *SimpleReq) getToken() string {

	if s.ReqType == ReqTypeForManger {
		return s.Mangersign()
	} else if s.ReqType == ReqTypeForUpload {
		return ""
	} else {
		return sumDownToken(s.ak, s.sk, s.URI)
	}
}

// Mangersign for manger token
func (s *SimpleReq) Mangersign() string {
	str := s.Method + " " + s.URI + "\nHost: " + s.Host + "\nContent-Type: " + s.ContentType
	if s.XQinius != nil {
		str += "\n" + s.XQinius.ToSignURLString()
	}
	str += "\n\n"
	if s.ContentType != "application/octet-stream" {
		if s.Param != nil {
			str += s.Param.ToString()
		}
	}
	secret := hmac.New(sha1.New, []byte(s.sk))
	secret.Write([]byte(str))
	sign := safeBase64Encode(secret.Sum(nil))
	return fmt.Sprintf("Qiniu %s:%s", s.ak, sign)
}

// setBody setting http request Reader
func (s *SimpleReq) setBody() io.Reader {
	if s.Param != nil {
		return ioutil.NopCloser(bytes.NewReader(s.jsonBytes()))
	}
	return nil
}

// jsonStr json encode and base64.URLEncoding tostring
func (s *SimpleReq) jsonBytes() []byte {
	b, _ := json.Marshal(s.Param)
	return b
}

// AKSK 七牛ak,sk
type AKSK struct {
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
}

// IQiniuConf 七牛aksk 接口
type IQiniuConf interface {
	GetAKSK() (ak string, sk string)
}

// // SetConf init config
// func (a *AKSK) SetConf(cfunc ConfParseFunc) {
// 	cfunc(a)
// }

// XQiniu hanlder key such like  "X-Qiniu-<key>"  request headers kv string
type XQiniu struct {
	data map[string]string
}

// Set kv key will be turn lower and key no have "X-Qiniu-"
func (x *XQiniu) Set(key string, val string) {
	if key == "" {
		return
	}
	if x.data == nil {
		x.data = map[string]string{}
	}
	key = strings.ToLower(key)
	x.data[key] = val
}

// ToSignURLString turn to sign string
func (x *XQiniu) ToSignURLString() string {
	keys := []string{}
	for k, _ := range x.data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	str := []string{}
	for _, k := range keys {
		str = append(str, fmt.Sprintf("%s: %s", k, x.data[k]))
	}
	return strings.Join(str, "\n")
}

// SetToReqHeaders setting header kv
func (x *XQiniu) SetToReqHeaders(req *http.Request) {
	for k, val := range x.data {
		req.Header.Set("X-Qiniu-"+k, val)
	}
}
