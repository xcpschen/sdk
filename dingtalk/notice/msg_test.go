package notice_test

import (
	"sdk/dingtalk/notice"
	"testing"
	"time"
)

var msg = `{"Code":"","Name":"","CorpId":"ding378105c44217be81","CorpSecret":"dbC6SykQDgvEm3PckDfPADMcGyOFQc-q_0dWRUzpTpDQ7iWPEIt3kpn09rZP4cVV","AppId":"dingoasi1kq71qjmu7pjw6","AppSecret":"ownpS3JOOiFVGMBKgE5lYWp-3CNpg0Ld8BMwWmMpP6YzpyYfD5eK9MyMPyporlzX","AgentId":"205149371"}`

func TestWm(t *testing.T) {
	w := new(notice.WebHook)
	w.URL = "https://oapi.dingtalk.com/robot/send?access_token=c72dee25f2446861475182e0eaa1c5a0ead09f7e7047501b0e853f84c04493e7"
	md := new(notice.Markdown)
	md.AddTitle("新工单", 4)
	md.AddSplitline()
	md.AddList(true, []string{
		"工单编号：6DKkfWV109000001",
		"提单人：xxxxxx",
		"工单类型:计算机-蓝屏",
	}...)
	md.AddItalicText(time.Now().Format("2006-01-02 15:04"))
	mb := &notice.ContentMarkDown{
		Title: "新工单",
		MD:    md,
	}

	w.Send(mb)
}
