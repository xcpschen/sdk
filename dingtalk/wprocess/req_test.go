package wprocess_test

import (
	"dingtalk/common"
	"dingtalk/wprocess"
	"fmt"
	"testing"
)

const (
	ac string = "bcae3eec338a38f381d9993a3f811cf4"
)

func TestReq(t *testing.T) {
	tc := wprocess.TaskCancel{
		Agentid:           "4321",
		ProcessInstanceID: "234567",
		ActiivityID:       "675432",
		// ActivityIDList:    []string{"324567", "67543"},
	}
	fmt.Println(common.JsonEncode(tc))
	t.Log("ok")
}

func TestCreateWorkRecord(t *testing.T) {
	w := &wprocess.CreateWorkRecord{
		ProcessCode:      "PROC-53876730-E950-4C25-8E46-82ACCB9E33AC",
		OriginatorUserID: "manager8409",
		Title:            "测试sql执行",
		FormComponentValues: []wprocess.ComponentValues{
			wprocess.ComponentValues{
				Name:  "申请事项",
				Value: "测试输入框",
			},
			wprocess.ComponentValues{
				Name:  "执行目标库",
				Value: "多行value",
			},
			wprocess.ComponentValues{
				Name:  "执行时间",
				Value: "2020-02-03",
			},
		},
		URL: "http://cmdb.vaiwan.com/api/apibpms/test",
	}
	p := wprocess.NewDTProcessWithAccessToken(ac)
	if err := p.Exec(wprocess.DtReq(w)); err != nil {
		t.Fatal(err)
	}
	//{"errcode":0,"result":{"process_instance_id":"1a3a08ab-4c9b-4261-8045-18ee8485002c"},"success":true,"request_id":"y5xz16yppbkk"}
	fmt.Println(common.JsonEncode(w.Resp))
}

func TestUpdateWorkRecord(t *testing.T) {
	u := &wprocess.TaskCreate{
		ProcessInstanceID: "441b2216-48f7-4c6f-9281-7f3c9635c749",
		Tasks: []wprocess.DtTask{
			wprocess.DtTask{
				UserID: "016610005731975677",
				URL:    "http://cmdb.vaiwan.com/api/apibpms/test?UserID=016610005731975677",
			},
		},
	}
	p := wprocess.NewDTProcessWithAccessToken(ac)
	if err := p.Exec(wprocess.DtReq(u)); err != nil {
		t.Fatal(err)
	}
	//{"errcode":0,"tasks":[{"task_id":63266761388,"userid":"manager8409"},{"task_id":63266761389,"userid":"manager8409"}],"request_id":"6q28gf3arst9"}
	fmt.Println(common.JsonEncode(u.Resp))
}

func TestSaveTpl(t *testing.T) {
	u := &wprocess.SaveTpl{}
	u.Name = "sql审批1"
	u.Description = "sql审批"
	u.FakeMode = true
	c1 := wprocess.Component{
		Name: "TextField",
		Props: wprocess.Prop{
			ID:       "TextField-1234567a",
			Label:    "申请事项",
			Required: true,
		},
	}
	c2 := wprocess.Component{
		Name: "TextareaField",
		Props: wprocess.Prop{
			ID:       "TextareaField-1234567b",
			Label:    "执行目标库",
			Required: true,
		},
	}
	c3 := wprocess.Component{
		Name: "TextField",
		Props: wprocess.Prop{
			ID:       "TextField-1234567c",
			Label:    "执行时间",
			Required: true,
		},
	}
	u.SetFcList(c1, c2, c3)
	data, _ := u.Data()
	fmt.Println(common.JsonEncode(data))
	ac := "9bba05b751323312ad46d015bc71486d"
	p := wprocess.NewDTProcessWithAccessToken(ac)
	if err := p.Exec(wprocess.DtReq(u)); err != nil {
		t.Fatal(err)
	}

	//{"errcode":0,"tasks":[{"task_id":63266761388,"userid":"manager8409"},{"task_id":63266761389,"userid":"manager8409"}],"request_id":"6q28gf3arst9"}
	fmt.Println(common.JsonEncode(u.Resp))
}

func TestTaskCancel(t *testing.T) {
	w := &wprocess.TaskCancel{
		ProcessInstanceID: "63532f41-d536-4b69-a007-4bbb1cdc13e6",
	}
	p := wprocess.NewDTProcessWithAccessToken(ac)
	if err := p.Exec(wprocess.DtReq(w)); err != nil {
		t.Fatal(err)
	}
	fmt.Println(common.JsonEncode(w.Resp))
}

func TestQueryTask(t *testing.T) {
	q := &wprocess.QueryTask{
		UserID: "127657696",
		Count:  50,
	}
	data, _ := q.Data()
	fmt.Println(common.JsonEncode(data))
	p := wprocess.NewDTProcessWithAccessToken(ac)
	if err := p.Exec(wprocess.DtReq(q)); err != nil {
		t.Fatal(err)
	}
	fmt.Println(common.JsonEncode(q.Resp))
}

func TestUpdateProcess(t *testing.T) {
	// a801ffad-8d63-4042-9b69-3b604d150e59
	// 6c20349d-4afe-4d49-826e-f23606cd709f
	w := &wprocess.UpdateBatchWorkRecord{
		Agentid: "205149371",
		Instances: []wprocess.UpdateInstances{
			wprocess.UpdateInstances{
				ProcessInstanceID: "28ff4c4f-7b1c-4460-872a-5f7ca23a4725",
				Status:            "TERMINATED",
				Result:            "agree",
			},
		},
	}
	data, _ := w.Data()
	fmt.Println(common.JsonEncode(data))
	p := wprocess.NewDTProcessWithAccessToken(ac)
	// p := wprocess.NewDTProcessWithAccessToken("2b122165dd9c3e9fa9212d1b376afce2")
	if err := p.Exec(wprocess.DtReq(w)); err != nil {
		t.Fatal(err)
	}
	fmt.Println(common.JsonEncode(w.Resp))
}

type FlowDTTmp struct {
	Dtconf DtConf `json:"dtconf"`
}
type DtConf struct {
	From        string        `json:"from"`
	ProcessCode string        `json:"process_code"`
	Url         string        `json:"url"`
	Fvs         []FlowDTTmpFv `json:"fv"`
}
type FlowDTTmpFv struct {
	Label  string   `json:"label"`
	Vals   []string `json:"vals"`
	Format string   `json:"format"`
}

func TestAd(t *testing.T) {
	// q := `{"dtconf":{"form":{"saveProcessRequest":{"description":"sql审批","fake_mode":true,"form_component_list":[{"component_name":"TextField","props":{"id":"TextField-1234567a","label":"申请事项","required":true,"choiceomitempty":0,"stat_field":{"id":"","label":null,"upper":false,"unit":""}}},{"component_name":"TextareaField","props":{"id":"TextareaField-1234567b","label":"执行目标库","required":true,"choiceomitempty":0,"stat_field":{"id":"","label":null,"upper":false,"unit":""}}},{"component_name":"TextField","props":{"id":"TextField-1234567c","label":"执行时间","required":true,"choiceomitempty":0,"stat_field":{"id":"","label":null,"upper":false,"unit":""}}}],"name":"sql审批1","process_code":""}},"process_code":"PROC-53876730-E950-4C25-8E46-82ACCB9E33AC","url":"http://cmdb.vaiwan.com/api/apibpms/test","fv":[{"label":"申请事项","vals":["title"],"format":""},{"label":"执行目标库","vals":["instanceName","dbName"],"format":"实例%q,数据库%q"},{"label":"执行时间","vals":["runTime"],"format":""}]}}`
	// data := &FlowDTTmp{}
	// if err := json.Unmarshal([]byte(q), data); err != nil {
	// 	t.Fatal(err)
	// }
	str := "234"
	str2 := 234

	fmt.Println(fmt.Sprintf("%v,%v", str, str2))
}
