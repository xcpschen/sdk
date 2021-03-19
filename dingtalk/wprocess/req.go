package wprocess

import (
	"dingtalk/common"
	"fmt"
)

type (
	SaveTplRespone struct {
		DTRespone
		Result interface{} `json:"result"`
	}
	SaveTpl struct {
		ProcessFormData
		Resp *SaveTplRespone `json:"-"`
	}
)

func (st *SaveTpl) Uri(accessToken string) string {
	return fmt.Sprintf("/topapi/process/save?access_token=%s", accessToken)
}

func (st *SaveTpl) Method() string {
	return "POST"
}
func (st *SaveTpl) SetFcList(cp ...Component) {
	if st.FormComponentList == nil {
		st.FormComponentList = []Component{}
	}
	st.FormComponentList = append(st.FormComponentList, cp...)
}
func (st *SaveTpl) Data() (val interface{}, err error) {
	if st.Name == "" {
		return nil, fmt.Errorf("创建或跟新模版，name必须填写")
	}
	if st.Description == "" {
		return nil, fmt.Errorf("创建或跟新模版，description必须填写")
	}
	if len(st.FormComponentList) <= 0 {
		return nil, fmt.Errorf("创建或跟新模版，FormComponentList必须至少存在一个控件")
	}

	for i, item := range st.FormComponentList {
		if item.Props.ID == "" {
			item.Props.ID = fmt.Sprintf("%s-%s", item.Name, common.RandNumStr(8))
		}

		switch n := item.Props.Label.(type) {
		case string:
			if n == "" {
				return nil, fmt.Errorf("第%d控件标签label必须设置", i)
			}
		case []string:
			if len(n) == 0 {
				return nil, fmt.Errorf("第%d控件标签label必须设置", i)
			}
		default:
			return nil, fmt.Errorf("第%d控件标签label类型必须是string,or string 数组", i)
		}
	}

	st.FakeMode = true
	data := map[string]interface{}{
		"saveProcessRequest": *st,
	}
	return data, nil
}

func (st *SaveTpl) SetRespone(msg []byte) error {
	st.Resp = &SaveTplRespone{}
	if err := common.JsonDecodeByte(msg, &st.Resp); err != nil {
		return err
	}
	if st.Resp.ErrCode == ReqSucess {
		st.ProcessCode = st.Resp.Result.(map[string]interface{})["process_code"].(string)
	} else {
		return fmt.Errorf(st.Resp.ErrMsg)
	}
	return nil
}

type DeleteTpl struct {
	Agentid         string `json:"agentid,int64,omitempty"`
	ProcessCode     string `json:"process_code"`
	CleanRuningTask bool   `json:"clean_running_task,omitempty"`
	//是否清理运行中的任务。true表示是
	Resp *DTRespone `json:"-"`
}

func (dt *DeleteTpl) Uri(accessToken string) string {
	return fmt.Sprintf("/topapi/process/delete?access_token=%s", accessToken)
}

func (dt *DeleteTpl) Method() string {
	return "POST"
}

func (dt *DeleteTpl) Data() (val interface{}, err error) {
	if dt.ProcessCode == "" {
		return nil, fmt.Errorf("删除模版，ProcessCode必须填写")
	}
	data := map[string]interface{}{
		"request": *dt,
	}
	return data, nil
}

func (dt *DeleteTpl) SetRespone(msg []byte) error {
	dt.Resp = &DTRespone{}
	if err := common.JsonDecodeByte(msg, dt.Resp); err != nil {
		return err
	}
	return nil
}

type (
	CwrRespone struct {
		DTRespone
		Result            map[string]interface{} `json:"result"`
		ProcessInstanceID string                 `json:"process_instance_id"`
	}
	CreateWorkRecord struct {
		Agentid             string            `json:"agentid,int64,omitempty"`
		ProcessCode         string            `json:"process_code"`
		OriginatorUserID    string            `json:"originator_user_id"`
		Title               string            `json:"title,omitempty"`
		FormComponentValues []ComponentValues `json:"form_component_values"`
		URL                 string            `json:"url,omitempty"`
		Resp                *CwrRespone       `json:"-"`
	}
)

func (cr *CreateWorkRecord) Uri(accessToken string) string {
	return fmt.Sprintf("/topapi/process/workrecord/create?access_token=%s", accessToken)
}

func (cr *CreateWorkRecord) Method() string {
	return "POST"
}

func (cr *CreateWorkRecord) Data() (val interface{}, err error) {
	if cr.ProcessCode == "" {
		return nil, fmt.Errorf("创建实例，ProcessCode必须填写")
	}
	if cr.OriginatorUserID == "" {
		return nil, fmt.Errorf("创建实例，OriginatorUserID 发起人id必须填写")
	}

	if err := cr.checkFormValue(); err != nil {
		return nil, err
	}
	if cr.URL == "" {
		return nil, fmt.Errorf("创建实例，实例跳转url 必须填写")
	}
	data := map[string]interface{}{
		"request": *cr,
	}
	return data, nil
}
func (cr *CreateWorkRecord) checkFormValue() error {
	if len(cr.FormComponentValues) <= 0 {
		return fmt.Errorf("创建实例，请设置控件值")
	}
	return nil
}

func (cr *CreateWorkRecord) SetRespone(msg []byte) error {
	cr.Resp = &CwrRespone{}
	if err := common.JsonDecodeByte(msg, cr.Resp); err != nil {
		return err
	}
	if cr.Resp.ErrCode == ReqSucess {

		cr.Resp.ProcessInstanceID = cr.Resp.Result["process_instance_id"].(string)
		return nil
	} else {
		return fmt.Errorf("【%d】%s", cr.Resp.ErrCode, cr.Resp.ErrMsg)
	}
}

type UpdateWorkRecord struct {
	Agentid           string     `json:"agentid,int64,omitempty"`
	ProcessInstanceID string     `json:"process_instance_id"`
	Status            string     `json:"status"`
	Result            string     `json:"resule"`
	Resp              *DTRespone `json:"-"`
}

func (uw *UpdateWorkRecord) Uri(accessToken string) string {
	return fmt.Sprintf("/topapi/process/workrecord/update?access_token=%s", accessToken)
}

func (uw *UpdateWorkRecord) Method() string {
	return "POST"
}

func (uw *UpdateWorkRecord) Data() (val interface{}, err error) {
	status := []string{ProcessCompleted, ProcessTerminated}
	result := []string{ProcessAgree, ProcessRefuse}
	if uw.ProcessInstanceID == "" {
		return nil, fmt.Errorf("更新实例状态，ProcessInstanceID必须填写")
	}
	if uw.Status == "" {
		return nil, fmt.Errorf("更新实例状态，status 实例状态必须设置")
	} else if !common.InStrArray(uw.Status, status) {
		return nil, fmt.Errorf("更新实例状态，status 实例状态必须设置%qor%q", ProcessCompleted, ProcessTerminated)
	}
	if uw.Result == "" || !common.InStrArray(uw.Result, result) {
		return nil, fmt.Errorf("更新实例状态，Result任务结果必须设置分为agree和refuse")
	}
	data := map[string]interface{}{
		"request": *uw,
	}
	return data, nil
}
func (uw *UpdateWorkRecord) SetRespone(msg []byte) error {
	uw.Resp = &DTRespone{}
	if err := common.JsonDecodeByte(msg, uw.Resp); err != nil {
		return err
	}
	return nil
}

type UpdateBatchWorkRecord struct {
	Agentid   string            `json:"agentid,int64"`
	Instances []UpdateInstances `json:"instances"`
	Resp      *DTRespone        `json:"-"`
}

type UpdateInstances struct {
	ProcessInstanceID string    `json:"process_instance_id"`
	Status            string    `json:"status"`
	Result            string    `json:"result"`
	Resp              DTRespone `json:"-"`
}

func (ubw *UpdateBatchWorkRecord) Uri(accessToken string) string {
	return fmt.Sprintf("/topapi/process/workrecord/batchupdate?access_token=%s", accessToken)
}

func (ubw *UpdateBatchWorkRecord) Method() string {
	return "POST"
}

func (ubw *UpdateBatchWorkRecord) Data() (val interface{}, err error) {
	status := []string{ProcessCompleted, ProcessTerminated}
	result := []string{ProcessAgree, ProcessRefuse}
	if ubw.Agentid == "" {
		return nil, fmt.Errorf("批量更新实例状态，Agentid必须填写")
	}
	if len(ubw.Instances) <= 0 {
		return nil, fmt.Errorf("批量更新实例状态，Instances 实例列表必须设置")
	} else {
		for i, item := range ubw.Instances {
			if item.ProcessInstanceID == "" {
				return nil, fmt.Errorf("批量更新实例状态，第%d个，ProcessInstanceID必须填写", i)
			}
			if item.Status == "" {
				return nil, fmt.Errorf("批量更新实例状态，第%d个，status 实例状态必须设置", i)
			} else if !common.InStrArray(item.Status, status) {
				return nil, fmt.Errorf("批量更新实例状态，第%d个，status 实例状态必须设置%qor%q", i, ProcessCompleted, ProcessTerminated)
			}
			if item.Result == "" || !common.InStrArray(item.Result, result) {
				return nil, fmt.Errorf("批量更新实例状态，第%d个，Result任务结果必须设置分为agree和refuse", i)
			}
		}
	}
	data := map[string]interface{}{
		"request": *ubw,
	}
	return data, nil
}
func (ubw *UpdateBatchWorkRecord) SetRespone(msg []byte) error {
	ubw.Resp = &DTRespone{}
	if err := common.JsonDecodeByte(msg, ubw.Resp); err != nil {
		return err
	}
	return nil
}

type (
	TaskCreate struct {
		Agentid           string     `json:"agentid,string,omitempty"`
		ProcessInstanceID string     `json:"process_instance_id"`
		ActiivityID       string     `json:"json:"activity_id"`
		Tasks             []DtTask   `json:"tasks"`
		Resp              *TCRespone `json:"-"`
	}
	DtTask struct {
		TaskID int64  `json:"task_id,omitempty"`
		UserID string `json:"userid,omitempty"`
		Status string `json:"status,omitempty"`
		Result string `json:"result,omitempty"`
		URL    string `json:"url,omitempty"`
	}
	TCRespone struct {
		DTRespone
		Tasks []DtTask `json:"tasks"`
	}
)

func (tc *TaskCreate) Uri(accessToken string) string {
	return fmt.Sprintf("/topapi/process/workrecord/task/create?access_token=%s", accessToken)
}

func (tc *TaskCreate) Method() string {
	return "POST"
}

func (tc *TaskCreate) Data() (val interface{}, err error) {
	if tc.ProcessInstanceID == "" {
		return nil, fmt.Errorf("创建代办事项，ProcessInstanceID必须填写")
	}
	// if tc.ActiivityID == "" {
	// 	return nil, fmt.Errorf("创建代办事项，ActiivityID必须填写")
	// }

	if len(tc.Tasks) <= 0 {
		return nil, fmt.Errorf("创建代办事项，Tasks必须设置")
	} else {
		for i, item := range tc.Tasks {
			if item.UserID == "" {
				return nil, fmt.Errorf("创建代办事项，第%d个task UserID必须填写", i)
			} else if item.URL == "" {
				return nil, fmt.Errorf("创建代办事项，第%d个task URL必须填写", i)
			}
		}
	}
	data := map[string]interface{}{
		"request": *tc,
	}
	return data, nil
}

func (tc *TaskCreate) SetRespone(msg []byte) error {
	tc.Resp = &TCRespone{}
	if err := common.JsonDecodeByte(msg, tc.Resp); err != nil {
		return err
	}
	return nil
}

type TaskUpdate struct {
	Agentid           string     `json:"agentid,int64,omitempty"`
	ProcessInstanceID string     `json:"process_instance_id"`
	Tasks             []DtTask   `json:"tasks"`
	Resp              *DTRespone `json:"-"`
}

func (tu *TaskUpdate) Uri(accessToken string) string {
	return fmt.Sprintf("/topapi/process/workrecord/task/update?access_token=%s", accessToken)
}

func (tu *TaskUpdate) Method() string {
	return "POST"
}

func (tu *TaskUpdate) Data() (val interface{}, err error) {
	if tu.ProcessInstanceID == "" {
		return nil, fmt.Errorf("更新待办状态，ProcessInstanceID必须填写")
	}

	if len(tu.Tasks) <= 0 {
		return nil, fmt.Errorf("更新待办状态，Tasks必须设置")
	} else {
		for i, item := range tu.Tasks {
			if item.TaskID == 0 {
				return nil, fmt.Errorf("更新待办状态，第%d个task TaskID必须填写", i)
			} else if item.Status == "" {
				return nil, fmt.Errorf("更新待办状态，第%d个taskStatus必须填写", i)
			} else if item.Result == "" {
				return nil, fmt.Errorf("更新待办状态，第%d个task Result必须填写", i)
			}
		}
	}
	data := map[string]interface{}{
		"request": *tu,
	}
	return data, nil
}
func (tu *TaskUpdate) SetRespone(msg []byte) error {
	tu.Resp = &DTRespone{}
	if err := common.JsonDecodeByte(msg, tu.Resp); err != nil {
		return err
	}
	return nil
}

type TaskCancel struct {
	Agentid           string     `json:"agentid,int64,omitempty"`
	ProcessInstanceID string     `json:"process_instance_id"`
	ActiivityID       string     `json:"activity_id,omitempty"`
	ActivityIDList    []string   `json:"activity_id_list,omitempty"`
	Resp              *DTRespone `json:"-"`
}

func (tc *TaskCancel) SetRespone(msg []byte) error {
	tc.Resp = &DTRespone{}
	if err := common.JsonDecodeByte(msg, tc.Resp); err != nil {
		return err
	}
	return nil
}

func (tc *TaskCancel) Uri(accessToken string) string {
	return fmt.Sprintf("/topapi/process/workrecord/taskgroup/cancel?access_token=%s", accessToken)
}

func (tc *TaskCancel) Method() string {
	return "POST"
}

func (tc *TaskCancel) Data() (val interface{}, err error) {
	if tc.ProcessInstanceID == "" {
		return nil, fmt.Errorf("批量取消待办，ProcessInstanceID必须填写")
	}
	data := map[string]interface{}{
		"request": *tc,
	}
	return data, nil
}

type (
	QueryTask struct {
		UserID string `json:"userid"`
		Offset int    `json:"offset"`
		//分页游标，从0开始，每次累加count的值
		Count int `json:"count",valid:"0-50"`
		//分页大小，最大50
		Status int `json:"status"`
		//0表示待处理，-1表示已经移除
		Resp *QTRespone `json:"-"`
	}
	QTRespone struct {
		DTRespone
		Result ResponeResult `json:"result"`
	}
	ResponeResult struct {
		HasMore bool        `json:"has_more"`
		List    ResponeList `json:"list"`
	}
	ResponeList struct {
		URL        string        `json:"url"`
		TaskID     string        `json:"taks_id"`
		InstanceID string        `json:"intance_id"`
		Title      string        `json:"title"`
		Forms      []ResponeForm `json:"forms"`
	}
	ResponeForm struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
)

func (qt *QueryTask) Uri(accessToken string) string {
	return fmt.Sprintf("/topapi/process/workrecord/task/query?access_token=%s", accessToken)
}

func (qt *QueryTask) Method() string {
	return "POST"
}

func (qt *QueryTask) Data() (val interface{}, err error) {
	if qt.UserID == "" {
		return nil, fmt.Errorf("查询待办列表,userid必须填写")
	}
	if qt.Count > 50 {
		return nil, fmt.Errorf("查询待办列表,count最大50")
	}
	if qt.Status != 0 {
		if qt.Status != -1 {
			return nil, fmt.Errorf("查询待办列表,status 0表示待处理，-1表示已经移除")
		}
	}
	return *qt, nil
}

func (qt *QueryTask) SetRespone(msg []byte) error {
	qt.Resp = &QTRespone{}
	fmt.Println(string(msg))
	if err := common.JsonDecodeByte(msg, qt.Resp); err != nil {
		return err
	}
	return nil
}
