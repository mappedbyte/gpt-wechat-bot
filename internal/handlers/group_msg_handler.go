package handlers

import (
	"github.com/eatmoreapple/openwechat"
	"math/rand"
	"strings"
	"time"
	"wechatBot/internal/gpt"
)

type GroupMessage struct {
	//消息
	msg *openwechat.Message
	//发送的用户
	sender *openwechat.User
	//自己
	self *openwechat.Self
	//群组
	group *openwechat.Group
}

func GroupMessageContextHandler() func(ctx *openwechat.MessageContext) {
	return func(ctx *openwechat.MessageContext) {
		groupMessage, err := NewGroupMessage(ctx.Message)
		if err != nil {
			return
		}
		err = groupMessage.handle()
		return
	}
}

func NewGroupMessage(message *openwechat.Message) (*GroupMessage, error) {
	sender, err := message.Sender()
	if err != nil {
		return nil, err
	}
	g := &openwechat.Group{User: sender}
	senderInGroup, err := message.SenderInGroup()
	if err != nil {
		return nil, err
	}
	return &GroupMessage{
		msg:    message,
		sender: senderInGroup,
		self:   sender.Self(),
		group:  g,
	}, nil
}

func (g *GroupMessage) handle() error {
	//如果是纯文本，使用ChatGPT进行回复
	if g.msg.IsText() && g.msg.IsAt() {
		return g.ReplyText()
	}
	return nil
}

func (g *GroupMessage) ReplyText() error {
	maxInt := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(2)
	time.Sleep(time.Duration(maxInt) * time.Second)
	ai := gpt.OpenAI{}
	responseText, err := ai.Chat(g.NewRequestText())
	if err != nil {

	}
	replaceText := "@" + g.self.NickName
	content := strings.ReplaceAll(g.msg.Content, replaceText, "")
	repeat := strings.Split(content, "\n")
	lens := len(repeat)
	if lens <= 50 {
		lens = 50
	}
	line := strings.Repeat("-", lens)
	atText := "@" + g.sender.NickName

	responseText = atText + "\u2005" + "\n" + content + "\n" + line + "\n" + responseText
	responseText = strings.Trim(responseText, "\n")
	_, err = g.msg.ReplyText(responseText)
	return err

}

func (g *GroupMessage) NewRequestText() []any {
	var requestList = make([]any, 0)
	requestList = append(requestList, gpt.Message{Role: "system",
		Content: "\nYou are ChatGPT, a large language model trained by OpenAI.\nKnowledge cutoff: 2021-09\nCurrent model: gpt-4\nCurrent time: 2024/1/12 19:18:07\nLatex inline: $x^2$ \nLatex block: $$e=mc^2$$\n\n",
	})
	replaceText := "@" + g.self.NickName
	content := strings.ReplaceAll(g.msg.Content, replaceText, "")
	requestList = append(requestList, gpt.Message{
		Role:    "user",
		Content: content,
	})
	return requestList
}
