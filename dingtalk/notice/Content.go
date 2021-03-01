package notice

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type IDict map[string]interface{}
type MsgBody interface {
	ToBody() IDict
	ToJsonStr() string
}

// 文本
type ContentText struct {
	Content string `json:"content"`
}

func (this *ContentText) GetTypeKey() string {
	return "text"
}
func (this *ContentText) ToBody() IDict {
	key := this.GetTypeKey()
	return IDict{
		"msgtype": key,
		key:       this,
	}
}
func (this *ContentText) ToJsonStr() string {
	b, _ := json.Marshal(this.ToBody())
	return string(b)
}

// 图片
type ContentImage struct {
	MediaId string `json:"media_id"` // 媒体文件id，可以通过媒体文件接口上传图片获取
}

func (this *ContentImage) GetTypeKey() string {
	return "image"
}
func (this *ContentImage) ToBody() IDict {
	key := this.GetTypeKey()
	return IDict{
		"msgtype": key,
		key:       this,
	}
}
func (this *ContentImage) ToJsonStr() string {
	b, _ := json.Marshal(this.ToBody())
	return string(b)
}

// 语音
type ContentVoice struct {
	MediaId  string `json:"media_id"` // 媒体文件id。2MB，播放长度不超过60s，AMR格式
	Duration string `json:"duration"` // 正整数，小于60，表示音频时长
}

func (this *ContentVoice) GetTypeKey() string {
	return "voice"
}
func (this *ContentVoice) ToBody() IDict {
	key := this.GetTypeKey()
	return IDict{
		"msgtype": key,
		key:       this,
	}
}
func (this *ContentVoice) ToJsonStr() string {
	b, _ := json.Marshal(this.ToBody())
	return string(b)
}

// 文件
type ContentFile struct {
	MediaId string `json:"media_id"` // 媒体文件id。引用的媒体文件最大10MB
}

func (this *ContentFile) GetTypeKey() string {
	return "file"
}
func (this *ContentFile) ToBody() IDict {
	key := this.GetTypeKey()
	return IDict{
		"msgtype": key,
		key:       this,
	}
}
func (this *ContentFile) ToJsonStr() string {
	b, _ := json.Marshal(this.ToBody())
	return string(b)
}

// 链接
type ContentLink struct {
	Title   string `json:"title"`      // 消息标题
	Content string `json:"text"`       // 消息内容, 过长只会展示部分
	Url     string `json:"messageUrl"` // 消息跳转地址
	PicUrl  string `json:"picUrl"`     // 图片地址
}

func (this *ContentLink) GetTypeKey() string {
	return "link"
}
func (this *ContentLink) ToBody() IDict {
	key := this.GetTypeKey()
	return IDict{
		"msgtype": key,
		key:       this,
	}
}
func (this *ContentLink) ToJsonStr() string {
	b, _ := json.Marshal(this.ToBody())
	return string(b)
}

// markdown
type ContentMarkdown struct {
	Title   string `json:"title"` // 首屏会话透出的展示内容
	Content string `json:"text"`  // 消息内容, markdown 格式
}

func (this *ContentMarkdown) GetTypeKey() string {
	return "markdown"
}
func (this *ContentMarkdown) ToJsonStr() string {
	b, _ := json.Marshal(this.ToBody())
	return string(b)
}
func (this *ContentMarkdown) ToBody() IDict {
	key := this.GetTypeKey()
	return IDict{
		"msgtype": key,
		key:       this,
	}
}

// 整体跳转
type ContentActionBtn struct {
	Title string `json:"title"`     // 标题
	Url   string `json:"actionURL"` // 链接地址
}
type ContentAction struct {
	Title           string              `json:"title"`                     // 首屏会话透出的展示内容
	MarkDownContent string              `json:"markdown"`                  // 消息内容, markdown格式的
	SingleTitle     string              `json:"single_title,omitempty"`    // 单个按钮的方案, 设置此值及SingleURL时 Btns 无效
	SingleURL       string              `json:"single_url,omitempty"`      // 点击SingleTitle按钮触发的URL
	BtnOrientation  string              `json:"btn_orientation,omitempty"` // 0按钮竖直排列，1按钮横向排列
	Btns            []*ContentActionBtn `json:"btn_json_list,omitempty"`   // 按钮列表
}

func (this *ContentAction) SetUrl(urlstr string) {
	urlstr = url.QueryEscape(urlstr)
	this.SingleURL = fmt.Sprintf("dingtalk://dingtalkclient/page/link?url=%s&pc_slide=true", urlstr)
}
func (this *ContentAction) ToBody() IDict {
	key := this.GetTypeKey()
	return IDict{
		"msgtype": key,
		key:       this,
	}
}
func (this *ContentAction) GetTypeKey() string {
	return "action_card"
}
func (this *ContentAction) ToJsonStr() string {
	b, _ := json.Marshal(this.ToBody())
	return string(b)
}

// Feed
type ContentFeedLink struct {
	Title  string `json:"title"`      // 标题
	Url    string `json:"messageURL"` // 链接
	PicUrl string `json:"picURL"`     // 图片地址
}

type ContentFeed struct {
	Links []*ContentFeedLink `json:"links"` // 链接列表
}

func (this *ContentFeed) GetTypeKey() string {
	return "feedCard"
}
func (this *ContentFeed) ToBody() IDict {
	key := this.GetTypeKey()
	return IDict{
		"msgtype": key,
		key:       this,
	}
}
func (this *ContentFeed) ToJsonStr() string {
	b, _ := json.Marshal(this.ToBody())
	return string(b)
}

// OA
type ContentOaHead struct {
	BgColor string `json:"bgcolor"` // 消息头部的背景颜色
	Text    string `json:"text"`    // 消息的头部标题 (向普通会话发送时有效，向企业会话发送时会被替换为微应用的名字)，长度限制为最多10个字符
}
type ContentOaForm struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
type ContentOaRich struct {
	Num  string `json:"num"`  // 单行富文本信息的数目
	Unit string `json:"unit"` //单行富文本信息的单位
}
type ContentOaBody struct {
	Title     string           `json:"title,omitempty"`      // 消息体的标题，建议50个字符以内
	Content   string           `json:"content,omitempty"`    // 消息体的内容，最多显示3行
	Image     string           `json:"image,omitempty"`      // 消息体中的图片，支持图片资源@mediaId
	FileCount string           `json:"file_count,omitempty"` // 自定义的附件数目。此数字仅供显示，钉钉不作验证
	Author    string           `json:"author,omitempty"`     // 自定义的作者名字
	Form      []*ContentOaForm `json:"form,omitempty"`       // 消息内容表达
	Rich      *ContentOaRich   `json:"rich,omitempty"`       // 单行富文本信息
}
type ContentOa struct {
	MsgUrl   string         `json:"message_url"`              // 消息点击链接地址，当发送消息为小程序时支持小程序跳转链接
	PcMsgUrl string         `json:"pc_message_url,omitempty"` // PC端点击消息时跳转到的地址
	Head     *ContentOaHead `json:"head"`                     // 消息头部内容
	Body     *ContentOaBody `json:"body"`                     // 消息体
}

func (this *ContentOa) GetTypeKey() string {
	return "oa"
}

func (this *ContentOa) ToBody() IDict {
	key := this.GetTypeKey()
	return IDict{
		"msgtype": key,
		key:       this,
	}
}
func (this *ContentOa) ToJsonStr() string {
	b, _ := json.Marshal(this.ToBody())
	return string(b)
}

type ContentMarkDown struct {
	Title     string
	MD        *Markdown
	AtMobiles []string
	IsAll     bool
}

func (this *ContentMarkDown) ToBody() IDict {
	return IDict{
		"msgtype": "markdown",
		"markdown": map[string]string{
			"title": this.Title,
			"text":  this.MD.GetText(),
		},
		"at": map[string]interface{}{
			"atMobiles": this.AtMobiles,
			"isAtAll":   this.IsAll,
		},
	}
}

func (this *ContentMarkDown) ToJsonStr() string {
	b, _ := json.Marshal(this.ToBody())
	return string(b)
}
