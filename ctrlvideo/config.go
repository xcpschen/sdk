package ctrlvideo

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"
)

// CtrlVideo ctrlvideo config struct
type CtrlVideo struct {
	CtrlVideoAPIHost string `json:"ctrlvideo_api_host"`
	AppID            string `json:"app_id"`
	AppKey           string `json:"app_key"`
}

const (
	// CtrlVideoAPIHost Api Host
	CtrlVideoAPIHost string = "https://apiivetest.ctrlvideo.com"
)

// Iconfig get request signtrue
type Iconfig interface {
	Sign(query string) string
	APIHost() string
}

// Sign return sign with query and appkey
func (c *CtrlVideo) Sign(query string) string {
	str := query + "&appkey=" + c.AppKey
	h := md5.New()
	h.Write([]byte(str))
	return strings.ToLower(hex.EncodeToString(h.Sum(nil)))
}

// APIHost return api host
func (c *CtrlVideo) APIHost() string {
	return c.CtrlVideoAPIHost
}

// GetProject get project from project id
func (c *CtrlVideo) GetProject(ProjectID string) (*Project, error) {
	p := &Project{
		CommRequest: NewCommRequest(c.AppID),
		ID:          ProjectID,
	}
	if err := p.Get(c); err != nil {
		return nil, err
	}
	return p, nil
}

// ReplaceProjectAssets replace project assets url
func (c *CtrlVideo) ReplaceProjectAssets(ProjectID string, list []ProjectAsset) error {
	if ProjectID == "" {
		return errors.New("请设置项目ID")
	}
	if len(list) == 0 {
		return errors.New("请设置替换资源")
	}
	b, _ := json.Marshal(list)
	rp := &ReplaceProjectAssest{
		CommRequest: NewCommRequest(c.AppID),
		ProjectID:   ProjectID,
		ReplaceList: string(b),
	}
	return rp.DoReplaceAssest(c)
}
