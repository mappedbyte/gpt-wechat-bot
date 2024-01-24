package handlers

import (
	"github.com/eatmoreapple/openwechat"
	"wechatBot/internal/global"
)

func NewHandlers() (msgFunc func(m *openwechat.Message), err error) {
	dispatcher := openwechat.NewMessageMatchDispatcher()
	//处理私信
	dispatcher.RegisterHandler(func(message *openwechat.Message) bool {
		return message.IsSendByFriend()
	}, UserMessageContextHandler())
	// 处理群消息
	dispatcher.RegisterHandler(func(message *openwechat.Message) bool {
		return message.IsSendByGroup()
	}, GroupMessageContextHandler())
	//处理好友添加
	dispatcher.RegisterHandler(func(message *openwechat.Message) bool {
		return message.IsFriendAdd()
	}, func(ctx *openwechat.MessageContext) {
		if global.ServerConfig.Chat.AutoPass {
			_, _ = ctx.Message.Agree("")
			return
		}
	})
	return dispatcher.AsMessageHandler(), nil
}
