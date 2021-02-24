package jh

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type FeedbackReq struct {
	AccessKey string `json:"accessKey"`
	ReqID     string `json:"requestId"`
	FType     string `json:"type"`
	LabelID   string `json:"labelId"`
}

var (
	//ShuMeiFeedbackFtypKeyMap 数美 type
	ShuMeiFeedbackFtypKeyMap = map[string]string{
		"error": "误杀",
		"miss":  "漏杀",
	}
	// ShuMeiFeedbackLabelMap 数美 labelId
	ShuMeiFeedbackLabelMap = map[string]string{
		"politics":   "涉政",
		"porn":       "色情",
		"sexy":       "性感",
		"ad":         "广告",
		"violence":   "暴恐",
		"ban":        "违禁",
		"logo":       "商企 logo",
		"qr":         "二维码",
		"socialFace": "社交人脸",
		"minor":      "未成年人",
		"star":       "公众人物",
	}
)

func (f *FeedbackReq) ToHTTPReq() (req *http.Request, err error) {
	var b []byte
	b, err = json.Marshal(f)
	fmt.Println(string(b))
	if err != nil {
		return
	}
	buf := bytes.NewReader(b)
	req, err = http.NewRequest("POST", FeedbackURL, buf)
	if err != nil {
		return
	}
	req.Header.Add("content-type", "application/json;charset=UTF-8")
	return
}
