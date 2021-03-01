package dingtalk

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"
)

type GetUserId struct {
	Code string
	Rsp  *GetUserIDRsp
}

type GetUserIDRsp struct {
	UserID   string `json:"userid"`
	SysLevel int    `json:"sys_level"`
	IsSys    bool   `json:"is_sys"`
	Respone
}

func (g *GetUserId) Uri(at string) string {
	return fmt.Sprintf("/user/getuserinfo?access_token=%s&code=%s", at, g.Code)
}

func (g *GetUserId) Method() string {
	return "GET"
}

func (g *GetUserId) Data() (interface{}, error) {
	if g.Code == "" {
		return nil, errors.New("code必须设置")
	}
	return nil, nil
}

func (g *GetUserId) SetRespone(msg []byte) error {
	g.Rsp = &GetUserIDRsp{}
	err := json.Unmarshal(msg, g.Rsp)
	if err != nil {
		return errors.New("钉钉结果解析错误")
	}
	if g.Rsp.ErrCode != ReqSucess {
		return errors.New(g.Rsp.ErrMsg)
	}
	return nil
}

type GetUser struct {
	UserId string
	Rsp    *GetUserRsp
}

type GetUserRsp struct {
	Respone
	UserID     string `json:"userid"`
	Unionid    string
	Name       string
	WorkPlace  string
	Mobile     string
	OrgEmail   string
	Email      string
	Department []int64
	Position   string
	Avatar     string
	Jobnumber  string
}

func (g *GetUser) Uri(at string) string {
	return fmt.Sprintf("/user/get?access_token=%s&userid=%s", at, g.UserId)
}

func (g *GetUser) Method() string {
	return "GET"
}

func (g *GetUser) Data() (interface{}, error) {
	if g.UserId == "" {
		return nil, errors.New("UserId必须设置")
	}
	return nil, nil
}

func (g *GetUser) SetRespone(msg []byte) error {
	g.Rsp = &GetUserRsp{}
	err := json.Unmarshal(msg, g.Rsp)
	if err != nil {
		return errors.New("钉钉结果解析错误")
	}
	if g.Rsp.ErrCode != ReqSucess {
		return errors.New(g.Rsp.ErrMsg)
	}
	return nil
}

type QrLogin struct {
	TmpAuthCode string
	TimeStamp   int64
	AccessKey   string
	AppSerect   string
	Signature   string
	Rsp         *QrRsp
}
type QrRsp struct {
	Respone
	UserInfo struct {
		Nick    string `json:"nick"`
		Openid  string `json:"openid"`
		Unionid string `json:"unionid"`
	} `json:"user_info"`
}

func (g *QrLogin) Uri(at string) string {
	u := url.Values{}
	u.Add("signature", g.Signature)
	u.Add("timestamp", fmt.Sprintf("%d", g.TimeStamp))
	u.Add("accessKey", g.AccessKey)

	return fmt.Sprintf("/sns/getuserinfo_bycode?%s", u.Encode())
}
func (g *QrLogin) DoSignature() {
	now := time.Now()
	g.TimeStamp = now.UnixNano() / 1e6
	mac := hmac.New(sha256.New, []byte(g.AppSerect))
	mac.Write([]byte(fmt.Sprintf("%d", g.TimeStamp)))
	bs := mac.Sum(nil)
	g.Signature = base64.StdEncoding.EncodeToString(bs)
}
func (g *QrLogin) Method() string {
	return "POST"
}

func (g *QrLogin) Data() (interface{}, error) {
	if g.TmpAuthCode == "" {
		return nil, errors.New("缺少临时码")
	}
	return map[string]string{
		"tmp_auth_code": g.TmpAuthCode,
	}, nil
}

func (g *QrLogin) SetRespone(msg []byte) error {
	g.Rsp = &QrRsp{}
	err := json.Unmarshal(msg, g.Rsp)
	if err != nil {
		return errors.New("钉钉结果解析错误")
	}
	if g.Rsp.ErrCode != ReqSucess {
		return errors.New(g.Rsp.ErrMsg)
	}
	return nil
}

type GetUseridByUnionid struct {
	Unionid string
	Rsp     struct {
		Respone
		ContactType int    `json:"contactType"`
		UserID      string `json:"userid"`
	}
}

func (g *GetUseridByUnionid) Uri(at string) string {
	u := url.Values{}
	u.Add("access_token", at)
	u.Add("unionid", g.Unionid)
	return fmt.Sprintf("/user/getUseridByUnionid?%s", u.Encode())
}
func (g *GetUseridByUnionid) Method() string {
	return "GET"
}

func (g *GetUseridByUnionid) Data() (interface{}, error) {
	if g.Unionid == "" {
		return nil, errors.New("缺少临时码unionid")
	}
	return nil, nil
}

func (g *GetUseridByUnionid) SetRespone(msg []byte) error {
	// g.Rsp = &QrRsp{}
	err := json.Unmarshal(msg, &g.Rsp)
	if err != nil {
		return errors.New("钉钉结果解析错误")
	}
	if g.Rsp.ErrCode != ReqSucess {
		return errors.New(g.Rsp.ErrMsg)
	}
	return nil
}
