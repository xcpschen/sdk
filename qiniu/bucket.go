package qiniu

import (
	"errors"
	"fmt"
	"time"
)

type Bucket struct {
	AKSK
	Name   string
	Region string
}

// Domain return domain list for bucket
func (b *Bucket) Domain() ([]DomainInfo, error) {
	req := &SimpleReq{
		ReqType: ReqTypeForManger,
		Host:    RS_API_HOST,
		URI:     "/v7/domain/list?tbl=" + b.Name,
		Method:  "GET",
	}
	req.SetAkSK2(b.AccessKey, b.SecretKey)
	client := new(Qclient)
	client.Timeout = 5 * time.Second

	if err := client.DoReq(req); err != nil {
		return nil, err
	}
	if req.Resp.Code != 200 {
		return nil, fmt.Errorf("Code:%d Err:%s", req.Resp.Code, req.Resp.Error)
	}
	if len(req.Resp.Body) == 0 {
		return []DomainInfo{}, nil
	}
	var data []DomainInfo
	err := req.Resp.ParseBody(&data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// FetchObject fetch object from other web store
func (b *Bucket) FetchObject(fo *FetchObject) (jobID string, err error) {
	// if b.Region == "" {
	// 	return "", errors.New("请设置区域")
	// }
	if fo.URL == "" {
		return "", errors.New("请设置抓起URL")
	}
	if fo.Bucket == "" {
		fo.Bucket = b.Name
	}
	if fo.Bucket == "" {
		return "", errors.New("请设置仓库")
	}
	req := &SimpleReq{
		ReqType:     ReqTypeForManger,
		Host:        RS_API_HOST,
		URI:         "/sisyphus/fetch",
		Method:      "POST",
		ContentType: "application/json",
		Param:       fo,
	}
	req.SetAkSK2(b.AccessKey, b.SecretKey)
	client := new(Qclient)
	client.Timeout = 5 * time.Second
	if err := client.DoReq(req); err != nil {
		return "", err
	}
	if req.Resp.Code != 200 {
		return "", fmt.Errorf("Code:%d Err:%s", req.Resp.Code, req.Resp.Error)
	}
	var rp FetchRspone
	err = req.Resp.ParseBody(&rp)
	if err != nil {
		return "", err
	}
	return rp.ID, nil
}
