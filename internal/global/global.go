package global

import (
	"net/http"
	"wechatBot/internal/config"
)

var (
	ServerConfig         *config.ServerConfig = &config.ServerConfig{}
	Client                                    = http.Client{}
	DeadlineExceededText                      = "请求GPT服务器超时[裂开]，请重新发送问题[旺柴]"
)
