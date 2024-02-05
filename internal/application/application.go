package application

import (
	"github.com/eatmoreapple/openwechat"
	"log/slog"
	"os"
	"wechatBot/internal/config"
	"wechatBot/internal/handlers"
	"wechatBot/internal/initialize"
	"wechatBot/internal/notify"
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

	bot.MessageErrorHandler = func(err error) error {
		slog.Info("bot.MessageErrorHandler", "exit", "微信Bot退出,开始执行推送逻辑")
		err = bot.CrashReason()
		if err != nil {
			slog.Error(err.Error())
		}
		pushPlus := notify.PushPlus{}
		notifyErr := pushPlus.SendNotify("微信Bot退出了,快去检查下!")
		if notifyErr != nil {
			slog.Error("push", "消息推送失败", notifyErr.Error())
		}
		return err
	}
	reloadStorage := openwechat.NewFileHotReloadStorage("storage.json")
	defer func() {
		_ = reloadStorage.Close()
	}()

	// 先尝试使用热登录
	err = bot.HotLogin(reloadStorage, openwechat.NewRetryLoginOption())
	if err != nil {
		// 如果热登录失败，再尝试使用扫码登录
		if err := bot.Login(); err != nil {
			slog.Error("application.Run", "errorMsg", "登录出现错误:"+err.Error())
			os.Exit(1)
		}
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
