package qiniu_test

import (
	"fmt"
	"sdk/qiniu"
	"testing"
)

const (
	vault_static_token_qiniu_access_key = ""
	vault_static_token_qiniu_secret_key = ""
)

// CallbackBody set
type CallbackBody struct {
	ProjectAssetID int64
	Magicvar       qiniu.MagicVar
}

var (
	aksk = qiniu.AKSK{
		AccessKey: vault_static_token_qiniu_access_key,
		SecretKey: vault_static_token_qiniu_secret_key,
	}
	bucket = &qiniu.Bucket{
		AKSK:   aksk,
		Name:   "dev-cdn-common",
		Region: "hd",
	}
	object = &qiniu.Object{
		AKSK:   aksk,
		Bucket: "dev-cdn-common",
	}
)

// Callbackbodytype get Callbackbodytype str
func (*CallbackBody) Callbackbodytype() string {
	return "application/x-www-form-urlencoded"
}

// Callbackurl return Callbackurl str
func (s *CallbackBody) Callbackurl(id int) string {
	return fmt.Sprintf("http://bsky.vaiwan.com/callback/qiniu/save/%d", id)
}

// Callbackbody str
func (s *CallbackBody) Callbackbody() string {
	s.Magicvar.Set("bucket")
	s.Magicvar.Set("key")
	s.Magicvar.Set("fname")
	s.Magicvar.Set("bodySha1")
	return s.Magicvar.ToString()
}

func TestFetchObj(t *testing.T) {
	mgv := new(CallbackBody)
	fobj := &qiniu.FetchObject{
		FileType:         qiniu.FileTypeNormal,
		IgnoreSameKey:    false,
		URL:              "https://apiivetest.ctrlvideo.com/openapi/ajax/get_ivideo_json?project_id=5118117409561481",
		Bucket:           bucket.Name,
		Key:              "dev-test/get_ivideo_json?project_id=5118117409561481",
		CallbackURL:      mgv.Callbackurl(1),
		Callbackbody:     mgv.Callbackbody(),
		Callbackbodytype: mgv.Callbackbodytype(),
	}
	id, err := bucket.FetchObject(fobj)
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println(id)
}

func TestDomain(t *testing.T) {
	domains, err := bucket.Domain()

	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println(domains)
}

func TestObjectPrivateURL(t *testing.T) {
	object.Key = "dev-test/0.jpg"
	url, err := object.GetPrivateURL()
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println(url)
	t.Log(url)
}

func TestObjectDisable(t *testing.T) {
	object.Key = "dev-test/0.jpg"
	err := object.Disable("dev-cdn-common", "dev-test/0.jpg")
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log("ok")
}

func TestObjectEnable(t *testing.T) {
	object.Key = "dev-test/0.jpg"
	err := object.Enable("dev-cdn-common", "dev-test/0.jpg")
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log("ok")
}
