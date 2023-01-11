package jh_test

import (
	"encoding/json"
	"fmt"
	"testing"
)

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
