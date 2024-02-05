package application

import (
	"github.com/eatmoreapple/openwechat"
	"log/slog"
	"os"
	"wechatBot/internal/config"
	"wechatBot/internal/handlers"
	"wechatBot/internal/initialize"
)

func Run() *openwechat.Bot {
	initialize.InitChooseConfig()
	initialize.InitProxy()
	initialize.InitDiscord()
	h, err := handlers.NewHandlers()
	if err != nil {
		slog.Error("application.Run", "errorMsg", "初始化消息处理器失败,"+err.Error())
		os.Exit(1)
	}
	bot := openwechat.DefaultBot(openwechat.Desktop)
	bot.UUIDCallback = config.CheckOs()
	bot.MessageHandler = h
	reloadStorage := openwechat.NewFileHotReloadStorage("storage.json")
	defer func() {
		_ = reloadStorage.Close()
	}()
	if err := bot.Login(); err != nil {
		slog.Error("application.Run", "errorMsg", "登录出现错误:"+err.Error())
		os.Exit(1)
	}
	user, err := bot.GetCurrentUser()
	if err != nil {
		slog.Error("application.Run", "errorMsg", "获取当前用户出现错误:"+err.Error())
		os.Exit(1)
	}
	slog.Info("application.Run", "当前登录用户", user.NickName)
	return bot
}

// 执行热登录
/*if err := bot.HotLogin(reloadStorage, openwechat.NewRetryLoginOption()); err != nil {
	slog.Error("登录出现错误:", err.Error())
	return bot, err
}*/
