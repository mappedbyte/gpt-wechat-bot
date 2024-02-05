package main

import (
	"log/slog"
	"wechatBot/internal/application"
	"wechatBot/internal/notify"
)

func main() {
	bot := application.Run()
	_ = bot.Block()
	slog.Info("application.Run", "MessageErrorHandler 微信Bot退出,开始执行推送逻辑")
	pushPlus := notify.PushPlus{}
	_ = pushPlus.SendNotify("微信Bot退出了,快去检查下!")
}
