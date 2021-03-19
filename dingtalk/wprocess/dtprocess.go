package wprocess

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type (
	DTProcess struct {
		accessToken string
		dtConf      IDtConf
		ac          *accessToken
	}
	ProcessFormData struct {
		Agentid                  string      `json:"agentid,omitempty"`
		Description              string      `json:"description,omitempty"`
		DisableFormEdit          string      `json:"disable_form_edit,omitempty"`
		DisableStopProcessButton bool        `json:"disable_stop_process_button,omitempty"`
		FakeMode                 bool        `json:"fake_mode"`
		FormComponentList        []Component `json:"form_component_list"`
		Hidden                   bool        `json:"hidden,omitempty"`
		Name                     string      `json:"name"`
		ProcessCode              string      `json:"process_code"`
		TemplateEditURL          string      `json:"template_edit_url,omitempty"`
	}
	Component struct {
		Name     string     `json:"component_name"`
		Props    Prop       `json:"props"`
		Children *Component `json:"children,omitempty"`
	}
	Prop struct {
		ID          string      `json:"id"`
		Label       interface{} `json:"label"`
		Required    bool        `json:"required"`
		NotPrint    int8        `json:"not_print,omitempty" valid:"0-1"`
		Placeholder string      `json:"placeholder,omitempty"`
		// 是否需要大写 默认是需要; 1:不需要大写, 空或者0:需要大写
		NotUpper int8 `json:"not_upper,omitempty" valid:"0-1"`
		// 数字组件/日期区间组件单位属性
		Unit string `json:"unit,omitempty"`
		//单选框或者多选框的选项多个参数请用","分隔
		Options string `json:"options,omitempty"`
		// 时间格式
		Format string `json:"format,omitempty"`
		//是否自动计算时长
		Duration bool `json:"duration,omitempty"`
		//内部联系人choice，1表示多选，0表示单选
		Choice int8 `json:"choiceomitempty" valid:"0-1"`
		// 说明文案的链接地址
		Link string `json:"link,omitempty"`
		//说明
		Content string `json:"content,omitempty"`
		// 增加明细动作名称
		ActionName string `json:"action_name,omitempty"`
		//需要计算总和的明细组件
		StatField StatField `json:"stat_field,omitempty"`
	}
	StatField struct {
		ID    string      `json:"id"`
		Lable interface{} `json:"label"`
		// 统计总和是否大写
		Upper bool `json:"upper"`
		//单位
		Unit string `json:"unit"`
	}
	ComponentValues struct {
		Name  string      `json:"name"`
		Value interface{} `json:"value,string"`
	}
	DTRespone struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}
)

//组件类型名称
const (
	//ComponentTextField 单输入框
	ComponentTextField string = "TextField"
	//ComponentTextareaField 多行输入框
	ComponentTextareaField string = "TextareaField"
	//ComponentMoneyField 金额输入框
	ComponentMoneyField string = "MoneyField"
	//ComponentNumberField 数字输入框
	ComponentNumberField string = "NumberField"
	//ComponentDDDateField 日期输入框
	ComponentDDDateField string = "DDDateField"
	//ComponentDDDateRangeField 日期区间
	ComponentDDDateRangeField string = "DDDateRangeField"
)

const (
	ProcessCompleted string = "COMPLETED"

	ProcessTerminated string = "TERMINATED"
	//TERMINATED表示实例被终止，此时钉钉会将当前实例所有正在进行中的任务状态，置为CANCELED。status为TERMINATED时，不需要传result。

	ProcessAgree  string = "agree"
	ProcessRefuse string = "refuse"

	dingtalkAPIHost string = "https://oapi.dingtalk.com"

	defaultReqTimeOut int = 60
	//默认5秒超时
)

const (
	ReqSucess int = 0
)

type DtReq interface {
	Uri(accessToken string) string
	Method() string
	Data() (val interface{}, err error)
	SetRespone(msg []byte) error
}

type IDtConf interface {
	GetCorpIDAndSecret() (string, string)
	GetAppIDAndSecret() (string, string)
	GetAgentID() string
}

func NewDTProcessWithAccessToken(accessToken string) *DTProcess {
	return &DTProcess{accessToken: accessToken}
}

func NewDTProcessWithIDtConf(c IDtConf) *DTProcess {
	return &DTProcess{dtConf: c}
}

//Exec 执行操作
//dtReq 必须实现DtReq 接口
func (dtp *DTProcess) Exec(dtReq DtReq) error {
	var url string
	if dtp.accessToken == "" {
		ac, err := dtp.getAccessToken()
		if err != nil {
			return err
		}
		url = fmt.Sprintf("%s%s", dingtalkAPIHost, dtReq.Uri(ac))
	} else {
		url = fmt.Sprintf("%s%s", dingtalkAPIHost, dtReq.Uri(dtp.accessToken))
	}
	method := dtReq.Method()
	param, err := dtReq.Data()
	if err != nil {
		return err
	}

	req, err := NewHttpRequest(url, method, param)
	if err != nil {
		return err
	}
	client := &http.Client{Timeout: time.Duration(defaultReqTimeOut) * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("网络请求错误:%s", err)
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("结果解析错误:%s", err)
	}
	if err := dtReq.SetRespone(result); err != nil {
		return err
	}
	return nil
}

func NewHttpRequest(url, method string, data interface{}) (req *http.Request, err error) {
	if data == nil {
		req, err = http.NewRequest(method, url, nil)
	} else {
		jsonStr, _ := json.Marshal(data)
		log.Println("【", method, "】", string(jsonStr))
		req, err = http.NewRequest(method, url, bytes.NewBuffer(jsonStr))
	}
	return
}

type accessToken struct {
	Value       string `json:"access_token"`
	ExpriesIn   int64  `json:"expires_in"`
	ExpriesDate int64
	DTRespone
}

func (dtp *DTProcess) getAccessToken() (string, error) {
	now := time.Now().Unix()
	if dtp.ac != nil && dtp.ac.ExpriesDate > now {
		return dtp.ac.Value, nil
	}
	appid, appsecret := dtp.dtConf.GetAppIDAndSecret()
	url := fmt.Sprintf("%s/gettoken?appkey=%s&appsecret=%s", dingtalkAPIHost, appid, appsecret)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {

		return "", fmt.Errorf("结果解析错误:%s", err)
	}
	dtp.ac = &accessToken{}
	if err := json.Unmarshal(result, dtp.ac); err != nil {
		return "", err
	}
	if dtp.ac.ErrCode != ReqSucess {
		return "", fmt.Errorf("%s", dtp.ac.ErrMsg)
	}
	dtp.ac.ExpriesDate = time.Now().Unix() + dtp.ac.ExpriesIn

	return dtp.ac.Value, nil
}
