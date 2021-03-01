package notice

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type WebHook struct {
	URL string
}

func (w *WebHook) Send(msg MsgBody) error {
	if w.URL == "" {
		return errors.New("请设置机器人地址")
	}
	body := msg.ToJsonStr()
	if body == "" {
		return errors.New("请设置机器人信息体")
	}
	client := &http.Client{}

	req, err := http.NewRequest("POST", w.URL, strings.NewReader(body))
	if err != nil {
		return errors.New("系统错误！！")
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("网络请求错误:%s", err)
	}
	_, err = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return fmt.Errorf("结果解析错误:%s", err)
	}

	return nil
}
