package global

import (
	"github.com/bwmarrin/discordgo"
	"net/http"
	"wechatBot/internal/config"
)

var (
	ServerConfig          *config.ServerConfig = &config.ServerConfig{}
	Client                                     = http.Client{}
	DiscordSession        *discordgo.Session
	DeadlineExceededText  = "请求GPT服务器超时[裂开]，请重新发送问题[旺柴]"
	DeadlineExceededImage = "https://raw.githubusercontent.com/oneAsiaPeople/gpt-wechat-bot/master/image/sorry.png"
)
