package notice

import (
	"encoding/json"
	"fmt"
	. "sdk/dingtalk"
	"strings"
)

type WorkMsg struct {
	Conf IDtConf
}

type MsgInfo struct {
	AgentID    string `json:"agent_id"`
	UserIDList string `json:"userid_list"`
	// DeptIDList string      `json:"dept_id_list"`
	ToAllUser bool        `json:"to_all_user"`
	Msg       interface{} `json:"msg"`

	Resp *MsgResp `json:"-"`
}
type MsgResp struct {
	Respone
	TaskID int64 `json:"task_id"`
}

func (this *MsgInfo) Uri(ac string) string {
	return fmt.Sprintf("/topapi/message/corpconversation/asyncsend_v2?access_token=%s", ac)
}
func (this *MsgInfo) Method() string {
	return "POST"
}
func (this *MsgInfo) Data() (interface{}, error) {
	if this.AgentID == "" {
		return nil, fmt.Errorf("Missing required arguments:agent_id")
	}
	if this.UserIDList == "" {
		return nil, fmt.Errorf("Missing required arguments:UserIDList")
	}
	if this.Msg == nil {
		return nil, fmt.Errorf("Missing required arguments:Msg ")
	}

	return this, nil
}
func (this *MsgInfo) SetRespone(msg []byte) error {
	this.Resp = &MsgResp{}

	if err := json.Unmarshal(msg, this.Resp); err != nil {
		return err
	}
	if this.Resp.ErrCode == ReqSucess {
		return nil
	}

	return fmt.Errorf("【%d】%s", this.Resp.ErrCode, this.Resp.ErrMsg)
}

func (this *WorkMsg) setReceiver(msg *MsgInfo, ids ...string) {
	var tmp []string
	if msg.UserIDList != "" {
		tmp = strings.Split(msg.UserIDList, ",")
	}
	for _, id := range ids {
		if !InStrArray(id, tmp) {
			tmp = append(tmp, id)
		}
	}
	msg.UserIDList = strings.Join(tmp, ",")
}

func (this *WorkMsg) Send(msg *MsgInfo) error {
	msg.AgentID = this.Conf.GetAgentID()
	client := NewClientWithIDtConf(this.Conf)
	if err := client.Exec(msg); err != nil {
		return err
	}

	if msg.Resp.ErrCode != ReqSucess {
		return fmt.Errorf("%d,%s", msg.Resp.ErrCode, msg.Resp.ErrMsg)
	}
	return nil
}
