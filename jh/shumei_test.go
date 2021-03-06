package jh_test

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestCheckOnt(t *testing.T) {
	// sm := jh.NewShuMeiWithConf(jh.ShuMeiConf{
	// 	AccessKey: "TnJ6uoPpUUeuvDMeh4vY",
	// 	AppID:     "default",
	// 	EventID:   "default",
	// 	TagType:   "POLITICS_VIOLENCE_BAN_PORN_MINOR_AD_LOGO",
	// 	TokenID:   "sXoTexBWs1dfyzt8eTev",
	// })
	// cf := &CFile{
	// 	Name:     "小清新.jpg",
	// 	URL:      "https://gimg2.baidu.com/image_search/src=http%3A%2F%2Fb-ssl.duitang.com%2Fuploads%2Fitem%2F201702%2F09%2F20170209222526_NmCsv.thumb.700_0.jpeg&refer=http%3A%2F%2Fb-ssl.duitang.com&app=2002&size=f9999,10000&q=a80&n=0&g=0n&fmt=jpeg?sec=1611991204&t=db0db45d64735ca202ed09dbbc887b3e",
	// 	FileHash: "3456543",
	// }
	// cf2 := &CFile{
	// 	Name:     "罗瑶.jpg",
	// 	URL:      "https://ss1.bdstatic.com/70cFuXSh_Q1YnxGkpoWK1HF6hhy/it/u=1071334779,862940228&fm=26&gp=0.jpg",
	// 	FileHash: "qwer22",
	// }
	// if err := sm.CheckMore([]jh.JHFile{cf, cf2}); err != nil {
	// 	t.Fatalf("err:%s", err.Error())
	// }
	// if err := sm.CheckOne(cf2); err != nil {
	// 	t.Fatalf("err:%s", err.Error())
	// }
	t.Log("ok")
}

type CFile struct {
	Name       string
	URL        string
	FileMime   string
	FileSize   int64
	FileHash   string
	jhSupplier string
	//鉴黄结构
	jhRspData interface{}
	//鉴黄是否出现异常，导致未鉴黄成功
	errorMsg string
}

// GetURL 返回七牛文件url
func (q *CFile) GetURL() string {
	return q.URL
}

// GetHash 返回文件hash
func (q *CFile) GetHash() string {
	return q.FileHash
}

// GetName 返回文件名
func (q *CFile) GetName() string {
	return q.Name
}

// GetMime 返回文件Mime
func (q *CFile) GetMime() string {
	return q.FileMime
}

// GetSize 返回文件大小
func (q *CFile) GetSize() int64 {
	return q.FileSize
}

// GetImgPixelSize 返回文件像素大小
func (q *CFile) GetImgPixelSize() int64 {
	return -1
}

// SaveJhResult 保存鉴黄
// param string driverName 鉴黄供应商
// param bool isPass 是否通过鉴黄
// param Object rspData 鉴黄结果
// param string errorMsg 异常错误
func (q *CFile) SaveJhResult(driverName string, isPass int, descri string, rspData string, errorMsg string) {
	q.jhSupplier = driverName
	q.jhRspData = rspData
	q.errorMsg = errorMsg
	b, _ := json.Marshal(rspData)
	fmt.Println(string(b))
}

func (q *CFile) ToLogString() string {
	data := map[string]interface{}{
		// "msg":        q.Callbackmsg,
		"jhSupplier": q.jhSupplier,
		"jhRspData":  q.jhRspData,
		"errorMsg":   q.errorMsg,
	}
	b, _ := json.Marshal(data)
	return string(b)
}
