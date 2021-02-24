package ctrlvideo

import (
	"fmt"
	"strings"
)

// Project info
type Project struct {
	CommRequest `url:"inherit"`
	ID          string         `url:"project_id"`
	List        []ProjectAsset `url:"--"`
}

// Get get project handle
func (P *Project) Get(c Iconfig) error {
	params := Param{}
	if err := StructToParam(P, &params); err != nil {
		return err
	}
	P.Sign = c.Sign(params.ToSignURLString())
	params.Set("sign", P.Sign)
	req, err := P.ToHTTPReq("GET", c.APIHost()+"/openapi/ajax/get_project_assets?"+params.ToSignURLString(), nil)
	fmt.Println(params.ToSignURLString())
	if err != nil {
		return err
	}
	client := &Hclient{}

	rsp, err := client.DoReq(req)
	if err != nil {
		return err
	}
	list := map[string][]ProjectAsset{}
	if err := rsp.Parse(&list); err != nil {
		return fmt.Errorf("结果解析失败,%v", rsp.Data)
	}
	P.List = list["list"]
	return nil
}

// ProjectAsset project asset item
type ProjectAsset struct {
	AssetID   string `json:"asset_id"`
	AssetType string `json:"asset_type"`
	URL       string `json:"url"`
}

// AssetName return asset name from url
func (p ProjectAsset) AssetName() string {
	arr := strings.Split(p.URL, "/")
	l := len(arr) - 1
	if l < 0 {
		return ""
	}
	if p.AssetType == "json" {
		a := strings.Split(arr[l], "=")
		al := len(a) - 1
		if al < 0 {
			return arr[l] + ".json"
		}
		return a[al] + ".json"
	}
	return arr[l]
}

// ReplaceProjectAssest replace Project assest body
type ReplaceProjectAssest struct {
	CommRequest `url:"inherit"`
	ProjectID   string `url:"project_id"`
	ReplaceList string `url:"replace_list"`
}

// DoReplaceAssest do replace assest handle
func (rp *ReplaceProjectAssest) DoReplaceAssest(c Iconfig) error {
	params := Param{}
	if err := StructToParam(rp, &params); err != nil {
		return err
	}
	rp.Sign = c.Sign(params.ToSignURLString())
	params.Set("sign", rp.Sign)

	reg, err := rp.ToHTTPReq("POST", c.APIHost()+"/openapi/ajax/replace_project_assets", strings.NewReader(params.ToSignURLString()))
	fmt.Println(c.APIHost() + "/openapi/ajax/replace_project_assets")
	fmt.Println(params.ToSignURLString())
	if err != nil {
		return err
	}
	client := &Hclient{}

	_, err = client.DoReq(reg)
	if err != nil {
		return err
	}
	return nil
}
