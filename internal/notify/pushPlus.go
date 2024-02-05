package notify

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"wechatBot/internal/global"
)

type PushPlus struct {
}

type PushPlusRequest struct {
	Token       string `json:"token"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	Template    string `json:"template"`
	Channel     string `json:"channel"`
	Webhook     string `json:"webhook"`
	CallbackUrl string `json:"callbackUrl"`
	Timestamp   string `json:"timestamp"`
}

func (p *PushPlus) SendNotify(msg string) error {
	plusRequest := PushPlusRequest{
		Token:   global.ServerConfig.PushPlus.Token,
		Title:   "微信机器人",
		Content: msg,
		Webhook: "wechat",
	}
	b, _ := json.Marshal(plusRequest)
	response, err := http.Post(global.ServerConfig.PushPlus.Url, "application/json", bytes.NewReader(b))
	if err != nil {
		return err
	}
	slog.Info("消息通知发送状态码:", "messageCode", response.StatusCode)
	return nil
}
