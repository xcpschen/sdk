package qiniu

import (
	"encoding/json"
)

// FetchObject fetch a Object from other web store
type FetchObject struct {
	URL              string `json:"url"`
	Bucket           string `json:"bucket"`
	Key              string `json:"key,omitempty"`
	Host             string `json:"host,omitempty"`
	Md5              string `json:"md5,omitempty"`
	Etag             string `json:"etag,omitempty"`
	CallbackURL      string `json:"callbackurl,omitempty"`
	Callbackbody     string `json:"callbackbody,omitempty"`
	Callbackbodytype string `json:"callbackbodytype,omitempty"`
	Callbackhost     string `json:"callbackhost,omitempty"`
	IgnoreSameKey    bool   `json:"ignore_same_key,omitempty"`
	FileType         int    `json:"file_type,omitempty"`
}

// FileType 类型
const (
	FileTypeNormal int = iota
	FileTypeLowAccess
	FileTypeStore
)

// ToParam to request params
func (f *FetchObject) ToParam() (data map[string]interface{}) {
	b, _ := json.Marshal(f)
	data = map[string]interface{}{}
	json.Unmarshal(b, &data)
	return
}

// ToString to string
func (f *FetchObject) ToString() string {
	b, _ := json.Marshal(f)
	return string(b)
}

// FetchRspone qiniu Fetch object Rspone
type FetchRspone struct {
	ID   string `json:"id"`
	Wait int64  `json:"wait"`
}


