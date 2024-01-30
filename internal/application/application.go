package application

import (
	"github.com/eatmoreapple/openwechat"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"os"
	"time"
	"wechatBot/internal/config"
	"wechatBot/internal/gin/middleware"
	"wechatBot/internal/gpt"
	"wechatBot/internal/handlers"
	"wechatBot/internal/initialize"
	"wechatBot/internal/notify"
)

func Run() *openwechat.Bot {
	initialize.InitConfig()
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
		slog.Info("application.Run", "MessageErrorHandler 微信Bot退出,开始执行推送逻辑")
		pushPlus := notify.PushPlus{}
		_ = pushPlus.SendNotify("微信Bot退出了,快去检查下!")
		return err
	}
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

type Prompt struct {
	Prompt string `json:"prompt"`
}

type DALLImage struct {
	Url string `json:"url"`
}

type DALLResponse struct {
	Created int64       `json:"created"`
	Data    []DALLImage `json:"data"`
}

func RunGin() {
	g := gin.Default()
	r := g.Group("/v1", middleware.Cors())
	r.POST("/images/generations", func(context *gin.Context) {
		prompt := Prompt{}
		dallImages := make([]DALLImage, 0)
		if err := context.ShouldBindJSON(&prompt); err == nil {
			images := gpt.ReplyImage(prompt.Prompt)

			for _, image := range images {
				dallImages = append(dallImages, DALLImage{
					Url: image,
				})
			}
			dallResponse := &DALLResponse{
				Created: time.Now().Unix(),
				Data:    dallImages,
			}
			context.JSON(http.StatusOK, dallResponse)
			return
		}

		dallResponse := &DALLResponse{
			Created: time.Now().Unix(),
			Data:    dallImages,
		}
		context.JSON(http.StatusOK, dallResponse)
	})

	g.Run(":12345")
}
