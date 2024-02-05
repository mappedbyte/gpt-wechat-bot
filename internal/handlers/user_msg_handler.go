package handlers

import (
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/patrickmn/go-cache"
	"log/slog"
	"math/rand"
	"net/http"
	"strings"
	"time"
	"wechatBot/internal/global"
	"wechatBot/internal/gpt"
	"wechatBot/internal/history"
)

var c = cache.New(60, time.Minute)

type UserMessage struct {
	//消息
	msg *openwechat.Message
	//发送的用户
	sender *openwechat.User
	//历史记录
	history history.History
}

func UserMessageContextHandler() func(ctx *openwechat.MessageContext) {
	return func(ctx *openwechat.MessageContext) {
		go func() {
			//获取消息
			userMessage, err := NewUserMessage(ctx.Message)
			if err != nil {
				slog.Error("init user message error")
				return
			}
			_ = userMessage.Handle()
		}()

	}
}
func NewUserMessage(message *openwechat.Message) (*UserMessage, error) {
	sender, err := message.Sender()
	if err != nil {
		return nil, err
	}
	userHistory := history.NewOpenAIUserHistory(sender, c)
	return &UserMessage{
		msg:     message,
		sender:  sender,
		history: userHistory,
	}, nil
}
func (u *UserMessage) Handle() error {
	if u.msg.IsText() {
		if strings.HasPrefix(u.msg.Content, "画") {
			return u.ReplyImage()
		}
		return u.ReplyText()
	}
	return nil
}

func (u *UserMessage) ReplyText() error {
	maxInt := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(2)
	time.Sleep(time.Duration(maxInt) * time.Second)
	ai := gpt.OpenAI{}
	responseText, err := ai.Chat(u.NewRequestText())
	fmt.Println(responseText)
	if err == nil {
		m := gpt.Message{
			Role:    "assistant",
			Content: u.msg.Content,
		}
		u.history.SetUserHistory(m)
	}
	_, err = u.msg.ReplyText(responseText)
	return err
}
func (u *UserMessage) NewRequestText() []any {
	h := u.history.GetHistory()
	m := gpt.Message{
		Role:    "user",
		Content: u.msg.Content,
	}
	h = append(h, m)
	u.history.SetUserHistory(m)
	return h
}

func (u *UserMessage) ReplyImage() error {
	content := u.msg.Content
	//content = strings.ReplaceAll(content, "画", "")

	images := gpt.ReplyImage(" " + content)
	//client := http.Client{}
	//client.Transport = global.Client.Transport
	for _, image := range images {
		if len(image) > 0 {
			request, _ := http.NewRequest("GET", image, nil)
			response, err := global.DiscordSession.Client.Do(request)
			if err != nil {
				continue
			}
			u.msg.ReplyImage(response.Body)
		}
	}

	return nil
}
