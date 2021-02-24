package qiniu

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

type Object struct {
	AKSK
	Bucket string
	Key    string
}

// Enable enable Object from qiniu stroe
func (o *Object) Enable(Bucket, key string) error {
	d := new(Dntry)
	d.Bucket = Bucket
	d.Key = key
	return o.setStatus(0, d.EncodedEntryURI())
}

// Disable disable Object from qiniu store
func (o *Object) Disable(Bucket, key string) error {
	d := new(Dntry)
	d.Bucket = Bucket
	d.Key = key
	return o.setStatus(1, d.EncodedEntryURI())
}

// setStatus modify qiniu Object status
// param status int 1 is forbidden,0 is activation
func (o *Object) setStatus(status int, encodedEntryURI string) error {
	req := &SimpleReq{
		ReqType: ReqTypeForManger,
		Host:    RS_BOX_HOST,
		URI:     fmt.Sprintf("/chstatus/%s/status/%d", encodedEntryURI, status),
		Method:  "POST",
	}
	req.SetAkSK2(o.AccessKey, o.SecretKey)
	client := new(Qclient)
	client.Timeout = 5 * time.Second

	if err := client.DoReq(req); err != nil {
		return err
	}
	if req.Resp.Code != 200 {
		return fmt.Errorf("Code:%d Err:%s", req.Resp.Code, req.Resp.Error)
	}
	return nil
}

// GetPrivateURL return a URL for private file
func (o *Object) GetPrivateURL(domains ...string) (string, error) {
	if len(domains) <= 0 {
		b := &Bucket{
			Name: o.Bucket,
		}
		b.AccessKey, b.SecretKey = o.AccessKey, o.SecretKey
		var (
			err        error
			domainInfo []DomainInfo
		)
		domainInfo, err = b.Domain()
		if err != nil {
			return "", err
		}
		domains = make([]string, len(domainInfo))
		for i, d := range domainInfo {
			domains[i] = d.Domain
		}
	}
	for _, d := range domains {
		uri := fmt.Sprintf("https://%s/%s?e=%d", d, o.Key, time.Now().Unix()+3600)
		uri = uri + "&token=" + sumDownToken(o.AccessKey, o.SecretKey, uri)
		r, _ := http.Head(uri)
		if r.StatusCode == 200 {
			return uri, nil
		}
	}
	return "", errors.New("获取不到正确的url")
}

// Move moving Object
func (o *Object) Move(destBucket string, destKey string, isForce bool) error {

	dest := new(Dntry)
	dest.Bucket = destBucket
	dest.Key = destKey

	src := new(Dntry)
	src.Bucket = o.Bucket
	src.Key = o.Key

	return o.reqTypeForManger(fmt.Sprintf("/move/%s/%s/force/%t", src.EncodedEntryURI(), dest.EncodedEntryURI(), isForce), "POST")
}

// Copy copy object
func (o *Object) Copy(destBucket string, destKey string, isForce bool) error {

	dest := new(Dntry)
	dest.Bucket = destBucket
	dest.Key = destKey

	src := new(Dntry)
	src.Bucket = o.Bucket
	src.Key = o.Key

	return o.reqTypeForManger(fmt.Sprintf("/copy/%s/%s/force/%t", src.EncodedEntryURI(), dest.EncodedEntryURI(), isForce), "POST")
}

func (o *Object) reqTypeForManger(URL, Method string) error {
	req := &SimpleReq{
		ReqType: ReqTypeForManger,
		Host:    RS_BOX_HOST,
		URI:     URL,
		Method:  Method,
	}
	req.SetAkSK2(o.AccessKey, o.SecretKey)
	client := new(Qclient)
	client.Timeout = 5 * time.Second

	if err := client.DoReq(req); err != nil {
		return err
	}
	if req.Resp.Code != 200 {
		return fmt.Errorf("Code:%d Err:%s", req.Resp.Code, req.Resp.Error)
	}
	return nil
}
