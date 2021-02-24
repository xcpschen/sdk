package jh

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
	// "github.com/astaxie/beego"
)

// ShuMeiConf shumei config
type ShuMeiConf struct {
	AccessKey string `json:"access_key"`
	AppID     string `json:"app_id"`
	EventID   string `json:"event_id"`
	TagType   string `json:"type"`
	TokenID   string `json:"tokenId"`
	//LabelBackList 检测标签被名单，如果命中该标签，且建议拒绝的无法通过
	LabelBackList []string `json:"label_back_list"`
}

const (
	//ShuMeiName 供应商名字
	ShuMeiName string = "数美"
	// ShuMeiConfCode config code
	ShuMeiConfCode string = "ShuMei"
	//ImgSizeLimit 图片大小限制
	ImgSizeLimit int64 = 10 * (1 << 20)
	//MinImgPixel 最小像素
	MinImgPixel int64 = 1 << 16
	//MaxCheckNum 最大检测数量
	MaxCheckNum int = 100

	defaultURL string = "sh"

	defaultAppID string = "default"

	defaultEventID string = "default"

	defaultTagType string = "POLITICS_VIOLENCE_BAN_PORN_MINOR_AD_LOGO"

	//SingleCheckTimeout  single file check http request timeout
	SingleCheckTimeout time.Duration = 5 * time.Second
	//BatchCheckTimeout  batch files check http request timeout
	BatchCheckTimeout time.Duration = 60 * time.Second
)

var (
	//SingleURLs single request map of region url
	SingleURLs = map[string]string{
		"bj":  "http://api-img-bj.fengkongcloud.com/image/v4",
		"sh":  "http://api-img-sh.fengkongcloud.com/image/v4",
		"sjp": "http://api-img-xjp.fengkongcloud.com/image/v4",
	}
	//BatchURLs batch request map of region url
	BatchURLs = map[string]string{
		//北京
		"bj": "http://api-img-bj.fengkongcloud.com/images/v4",
		//上海
		"sh": "http://api-img-sh.fengkongcloud.com/images/v4",
		//新加坡
		"sjp": "http://api-img-xjp.fengkongcloud.com/images/v4",
	}
	//FeedbackURL 图片纠错接口定义地址
	FeedbackURL = "http://api-web.fengkongcloud.com/image/feedback/v2"
	//Events 事件标识符
	Events = []string{
		"default",
		"headImage",
		"album",
		"dynamic",
		"article",
		"comment",
		"roomCover",
		"groupMessage",
		"message",
		"product",
	}
	//TagTypes 检测的风险类型
	TagTypes = []string{
		"POLITCS",
		"VIOLENCE",
		"BAN",
		"PORN",
		"MINOR",
		"AD",
		"LOGO",
		"STAR",
		"COR",
	}
)

// 请求状态码
const (
	BodyCodeSuccess            int = 1100
	BodyCodeQPSTimeout         int = 1101
	BodyCodeOverQPS            int = 1901
	BodyCodeIllegalParam       int = 1902
	BodyCodeServerFailure      int = 1903
	BodyCodeImgDownloadFailure int = 1911
	BodyCodeNoPrivilege        int = 9101
)

// 外置建议
const (
	RiskLevelPass   string = "PASS"
	RiskLevelReview string = "REVIEW"
	RiskLevelReject string = "REJECT"
)

// ShuMei 数美鉴黄
type ShuMei struct {
	Conf ShuMeiConf
}

// NewShuMeiWithConf new a instance of shumei with ShuMeiConf struct
func NewShuMeiWithConf(c ShuMeiConf) *ShuMei {
	return &ShuMei{
		Conf: c,
	}
}

// CheckOne 单个文件检测
func (sm *ShuMei) CheckOne(fs JHFile) error {
	r := new(ShuMeiReq)
	r.AccessKey = sm.Conf.AccessKey
	r.AppID = sm.Conf.AppID
	r.EventID = sm.Conf.EventID
	r.TagType = sm.Conf.TagType
	if err := r.setSingleReqData(sm.Conf.TokenID, fs); err != nil {
		return err
	}

	client := new(Hclient)

	client.Timeout = SingleCheckTimeout
	b, err := client.DoReq(r)
	if err != nil {
		return err
	}
	var rsp ShuMeiDetailRespone
	if err := JsonDecode(b, &rsp); err != nil {
		// beego.Error("shumei CheckOne %s", err.Error())
		// beego.Error(string(b))
		return err
	}

	if rsp.Code != BodyCodeSuccess {
		// beego.Error("shumeil CheckOne Error:%s", string(b))
		fs.SaveJhResult("ShuMeiName", JhIgnore, rsp.RequestID, "", JsonEncode(rsp), nil, fmt.Sprintf("code:%d,msg:%s", rsp.Code, rsp.Message))
		return nil
	}
	rsp.handleFsResult(fs, sm.Conf.LabelBackList)

	return nil
}

// CheckMore 多个文件检测
func (sm *ShuMei) CheckMore(fs []JHFile) error {
	r := new(ShuMeiReq)
	r.AccessKey = sm.Conf.AccessKey
	r.AppID = sm.Conf.AppID
	r.EventID = sm.Conf.EventID
	r.TagType = sm.Conf.TagType
	if err := r.setBatchReqData(sm.Conf.TokenID, fs); err != nil {
		return err
	}

	client := new(Hclient)
	client.Timeout = BatchCheckTimeout
	b, err := client.DoReq(r)
	if err != nil {
		return err
	}
	var rsp BatchShuMeiRespone
	if err := JsonDecode(b, &rsp); err != nil {
		// beego.Error("shumei CheckMore %s", err.Error())
		// beego.Error(string(b))
		return err
	}
	if rsp.Code != BodyCodeSuccess {
		// beego.Error("shumeil CheckMore Error:", string(b))
		for _, f := range fs {
			f.SaveJhResult("ShuMeiName", JhIgnore, rsp.RequestID, "", JsonEncode(rsp), nil, fmt.Sprintf("code:%d,msg:%s", rsp.Code, rsp.Message))
		}
		return nil
	}
	for _, img := range rsp.Imgs {
		key, _ := strconv.Atoi(img.BtID)
		f := fs[key]
		if img.Code != BodyCodeSuccess {
			f.SaveJhResult("ShuMeiName", JhIgnore, img.RequestID, "", JsonEncode(img), img.GetRikeTable(), fmt.Sprintf("code:%d,msg:%s", img.Code, img.Message))
		} else {
			img.handleFsResult(f, sm.Conf.LabelBackList)
		}
	}

	return nil
}

// Feedback shumei feddback miss
func (sm *ShuMei) Feedback(requestID string, ftype string, labelID ...string) error {
	if _, ok := ShuMeiFeedbackFtypKeyMap[ftype]; !ok {
		return fmt.Errorf("纠错类型错误")
	}
	f := new(FeedbackReq)
	f.AccessKey = sm.Conf.AccessKey
	f.ReqID = requestID
	f.FType = ftype
	if len(labelID) > 0 {
		f.LabelID = labelID[0]
	}
	client := new(Hclient)
	client.Timeout = 1 * time.Second

	b, err := client.DoReq(f)
	if err != nil {
		return err
	}
	var rsp BaseShuMeiRespone
	if err := JsonDecode(b, &rsp); err != nil {
		// beego.Error("shumei CheckMore %s", err.Error())
		// beego.Error(string(b))
		return err
	}
	if rsp.Code != BodyCodeSuccess {
		return errors.New(rsp.Message)
	}
	return nil
}

// ShuMeiReq http request data struct for shumei
type ShuMeiReq struct {
	AccessKey string      `json:"accessKey"`
	AppID     string      `json:"appId"`
	EventID   string      `json:"eventId"`
	TagType   string      `json:"type"`
	Data      interface{} `json:"data"`

	isSingle bool
}

func (r *ShuMeiReq) setSingleReqData(tokenID string, fs JHFile) error {
	imgURL := fs.GetURL()
	if imgURL == "" {
		// beego.Error(fs.GetName(), "获取不到图片地址")
		fs.SaveJhResult("ShuMeiName", JhIgnore, "", "", "", nil, "获取不到图片地址")
		return errors.New("没有可以执行的图片")
	}
	data := SimpleReqData{
		Img: fs.GetURL(),
	}
	data.TokenID = tokenID
	r.Data = data
	r.isSingle = true
	return nil
}

func (r *ShuMeiReq) setBatchReqData(tokenID string, fs []JHFile) error {
	imgs := []BatchImgs{}

	for i, f := range fs {
		imgURL := f.GetURL()
		if imgURL == "" {
			// beego.Error(f.GetName(), "获取不到图片地址")
			f.SaveJhResult("ShuMeiName", JhIgnore, "", "", "", nil, "获取不到图片地址")
			continue
		}
		imgs = append(imgs, BatchImgs{
			Img:  f.GetURL(),
			BtID: fmt.Sprintf("%d", i),
		})
	}
	if len(imgs) == 0 {
		return errors.New("没有可以执行的图片")
	}
	data := BatchReqData{
		Imgs: imgs,
	}
	data.TokenID = tokenID
	r.isSingle = false
	r.Data = data
	return nil
}

// ToHTTPReq 设置http请求体
func (r *ShuMeiReq) ToHTTPReq() (req *http.Request, err error) {
	var b []byte
	b, err = json.Marshal(r)
	fmt.Println(string(b))
	if err != nil {
		return
	}
	buf := bytes.NewReader(b)
	if r.isSingle {
		req, err = http.NewRequest("POST", SingleURLs["sh"], buf)
		if err != nil {
			return
		}
	} else {
		req, err = http.NewRequest("POST", BatchURLs["sh"], buf)
		if err != nil {
			return
		}
	}
	req.Header.Add("content-type", "application/json;charset=UTF-8")
	return
}

// SemoReqData 请求体data共同部分
type SemoReqData struct {
	TokenID  string       `json:"tokenId"`
	IP       string       `json:"ip,omitempty"`
	DeviceID string       `json:"deviceId,omitempty"`
	MaxFrame int          `json:"maxFrame,omitempty"`
	Interval int          `json:"interval"`
	Room     string       `json:"room,omitempty"`
	Extra    ReqDataExtra `json:"extra,omitempty"`
}

// ReqDataExtra 透传字段
type ReqDataExtra struct {
	//透传字段，Json
	PassThrought map[string]interface{} `json:"passThrough"`
}

// Add setting extra data with key and val
func (pt *ReqDataExtra) Add(key string, val interface{}) {
	if pt.PassThrought == nil {
		pt.PassThrought = map[string]interface{}{}
	}
	pt.PassThrought[key] = val
}

// ToString ReqDataExtra to string
func (pt *ReqDataExtra) ToString() string {
	b, _ := json.Marshal(pt.PassThrought)
	return string(b)
}

// SimpleReqData 单个文件检测请求结构
type SimpleReqData struct {
	SemoReqData
	Img string `json:"img"`
}

// BatchReqData batch file request struct
type BatchReqData struct {
	SemoReqData
	Imgs []BatchImgs `json:"imgs"`
}

// BatchImgs batch request imgs struct
type BatchImgs struct {
	Img  string `json:"img"`
	BtID string `json:"btI'"`
}

// BaseShuMeiRespone 数美返回体请求基本头部
type BaseShuMeiRespone struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"requestId"`
}

// ShuMeiDetailRespone 数美返回体审核文件基本信息
type ShuMeiDetailRespone struct {
	BaseShuMeiRespone
	// BtID 文件唯一标示
	BtID string `json:"btId"`
	//RiskLevel 外置建议，值可能是
	// PASS：正常，建议直接放行
	// REVIEW：可疑，建议人工审核
	// REJECT：违规，建议直接拦截
	RiskLevel string `json:"riskLevel"`
	// RiskLabel1 当riskLevel为PASS时返回normal
	RiskLabel1 string `json:"riskLabel1"`
	// RiskLabel2 当riskLevel为PASS时为空
	RiskLabel2 string `json:"riskLabel2"`
	// RiskLabel3 当riskLevel为PASS时为空
	RiskLabel3 string `json:"RiskLabel3"`
	// RiskDescription 风险描述
	RiskDescription string `json:"riskDescriotion"`
	// RiskDetail 风险详情
	RiskDetail RiskDetail `json:"riskDetail,omitempty"`
	// AuxInfo 其他辅助信息
	AuxInfo AuxInfoObj `json:"auxInfo,omitempty"`
	// AllLabels 返回命中的所有风险标签以及详情信息
	AllLabels []Labels `json:"allLabels,omitempty"`

	TokenLabels TokenLabel `json:"tokenLabels,omitempty"`
}

func (rsp ShuMeiDetailRespone) handleFsResult(fs JHFile, LabelBackList []string) {
	if rsp.RiskLevel == RiskLevelPass {
		fs.SaveJhResult("ShuMeiName", JhPass, rsp.RequestID, rsp.RiskDescription, JsonEncode(rsp), rsp.GetRikeTable(), "")
	} else {
		descri := rsp.RiskDescription
		if descri == "" {
			descri = rsp.RiskLabel1 + "|" + rsp.RiskLabel2 + "|" + rsp.RiskLabel3
		}
		b := InArr(rsp.RiskLabel1, LabelBackList)
		if b && rsp.RiskLevel == RiskLevelReject {
			//触发删除
			fs.SaveJhResult("ShuMeiName", JhReject, rsp.RequestID, descri, JsonEncode(rsp), rsp.GetRikeTable(), "")
		} else {
			//保留审核
			fs.SaveJhResult("ShuMeiName", JhReview, rsp.RequestID, descri, JsonEncode(rsp), rsp.GetRikeTable(), "")
		}
	}
}
func (rsp ShuMeiDetailRespone) GetRikeTable() []RiskLabel {
	data := make([]RiskLabel, len(rsp.AllLabels))
	for i, l := range rsp.AllLabels {
		data[i] = RiskLabel{
			Label: l.RiskLabel1,
			Score: l.Probability,
		}
	}
	return data
}

// TokenLabel d
type TokenLabel struct {
	UGCAccountRisk map[string]interface{} `json:"UGC_account_risk"`
}

// RiskDetail 风险详情
type RiskDetail struct {
	Faces   []RiskDetailFacesOrObjects `json:"faces,omitempty"`
	Objects []RiskDetailFacesOrObjects `json:"objects,omitempty"`
	OcrText RiskDetailOcrText          `json:"ocrText,omitempty"`
}
type RiskDetailFacesOrObjects struct {
	Name        string  `json:"name,omitempty"`
	Location    []int   `json:"location,omitempty"`
	Probability float64 `json:"probability,omitempty"`
}

type RiskDetailOcrText struct {
	Text         string            `json:"text,omitempty"`
	MatchedLists []OcrTextMl       `json:"matchedLists,omitempty"`
	RiskSegments []OcrTextSegments `json:"riskSegments,omitempty"`
}

type OcrTextMl struct {
	Name  string         `json:"name"`
	Words OcrTextMlWords `json:"words"`
}
type OcrTextMlWords struct {
	Word     string `json:"word"`
	Position []int  `json:"position,omitempty"`
}

type OcrTextSegments struct {
	Segment  string `json:"segment,omitempty"`
	Position []int  `json:"position,omitempty"`
}

type AuxInfoObj struct {
	Segments    int `json:"segments"`
	ErrorCode   int `json:"errorCode"`
	TypeVersion struct {
		POLITICS string `json:"POLITICS"`
		VIOLENCE string `json:"VIOLENCE"`
		BAN      string `json:"BAN"`
		PORN     string `json:"PORN"`
		MINOR    string `json:"MINOR"`
		AD       string `json:"AD"`
		SPAM     string `json:"SPAM"`
		LOGO     string `json:"LOGO"`
		STAR     string `json:"STAR"`
		OCR      string `json:"OCR"`
	} `json:"typeVersion,omitempty"`
	PassThrough map[string]interface{} `json:"passThrough,omitempty"`
}

// Labels matchs Labels detail
type Labels struct {
	RiskLabel1      string     `json:"riskLabel1"`
	RiskLabel2      string     `json:"riskLabel2"`
	RiskLabel3      string     `json:"riskLabel3"`
	RiskDescription string     `json:"riskDescription"`
	Probability     float64    `json:"probability"`
	RiskDetail      RiskDetail `json:"riskDetail,omitempty"`
}

// BatchShuMeiRespone batch requets respone for shumei
type BatchShuMeiRespone struct {
	BaseShuMeiRespone

	Imgs []ShuMeiDetailRespone `json:"imgs,omitempty"`

	AuxInfo AuxInfoObj `json:"auxInfo,omitempty"`
}
