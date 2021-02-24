package ctrlvideo_test

import (
	"fmt"
	"sdk/ctrlvideo"
	"testing"
	"time"
)

// ProjectAssets project asset item
type ProjectAssets struct {
	CommRequest `url:"inherit"`
	AssetID     string `url:"accest_id"`
	AssetType   string `url:"asset_type"`
	URL         string `url:"url"`
	NF          string `url:"--"`
	NFN         int    `url:"nfn"`
}

// CommRequest 公共请求参数
type CommRequest struct {
	TimeStamp int64  `url:"timestamp"`
	AppID     string `url:"appid"`
	Nonce     string `url:"nonce"`
	Sign      string `url:"sign"`
}
type C struct {
	AcName string
}

func TestStructToURLValues(t *testing.T) {
	data := ctrlvideo.Param{}
	p := &ProjectAssets{
		AssetID:   "2345",
		AssetType: "fdsfgfds",
		URL:       "http://www.baidu.com",
		NF:        "nf",
		NFN:       100,
	}
	p.TimeStamp = time.Now().Unix()
	p.AppID = "app_id"
	p.Nonce = ctrlvideo.RandString(16)
	err := ctrlvideo.StructToParam(p, &data)
	if err != nil {
		t.Error(err)
	}
	str := data.ToUrlEncodeString()
	fmt.Println(data.ToSignURLString())
	fmt.Println(str)
}

func TestGetProject(t *testing.T) {
	cv := &ctrlvideo.CtrlVideo{
		CtrlVideoAPIHost: "https://apiivetest.ctrlvideo.com",
		AppID:            "",
		AppKey:           "",
	}
	p, _ := cv.GetProject("5118117409561481")
	fmt.Println(p.List)
}

func TestReplaceAssets(t *testing.T) {
	cv := &ctrlvideo.CtrlVideo{
		CtrlVideoAPIHost: "https://apiivetest.ctrlvideo.com",
		AppID:            "",
		AppKey:           "",
	}
	asesst := []ctrlvideo.ProjectAsset{
		ctrlvideo.ProjectAsset{
			AssetID:   "3799498506167636",
			AssetType: "image",
			// https://res-1300249927.file.myqcloud.com/media/79/79/image/3799498506167636/source.gif
			URL: "https://res-1300249927.file.myqcloud.com/media/79/79/image/3799498506167636/source.gif",
		},
		ctrlvideo.ProjectAsset{
			AssetID:   "video_5118117409561481",
			AssetType: "video",
			//https://apiivetest.ctrlvideo.com/data/skyfall/media/18/118/release/5118117409561481/5118117409561481_20201204111902.mp4
			URL: "https://apiivetest.ctrlvideo.com/data/skyfall/media/18/118/release/5118117409561481/5118117409561481_20201204111902.mp4",
		},
		ctrlvideo.ProjectAsset{
			AssetID:   "json_5118117409561481",
			AssetType: "json",
			// https://apiivetest.ctrlvideo.com/openapi/ajax/get_ivideo_json?project_id=5118117409561481
			URL: "https://apiivetest.ctrlvideo.com/openapi/ajax/get_ivideo_json?project_id=5118117409561481",
		},
	}
	if err := cv.ReplaceProjectAssets("5118117409561481", asesst); err != nil {
		t.Fatal(err.Error())
	}
	p, _ := cv.GetProject("5118117409561481")
	fmt.Println(p.List)
}
