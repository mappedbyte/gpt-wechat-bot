package initialize

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"sync"
	"wechatBot/internal/config"
	"wechatBot/internal/global"
	"wechatBot/internal/gpt"
)

var configLock sync.Mutex

func InitProxy() {
	if global.ServerConfig.Chat.Proxy {
		if global.ServerConfig.Chat.ProxyUrl != "" {
			proxy, err := url.Parse(global.ServerConfig.Chat.ProxyUrl)
			if err != nil {
				slog.Error("initialize", "InitProxy error", fmt.Errorf("初始化代理失败,请检查配置文件和代理服务器状态: %s \n", err.Error()))
				os.Exit(1)
			}
			transport := &http.Transport{Proxy: http.ProxyURL(proxy)}
			global.Client.Transport = transport
		}
	}
}

func InitChooseConfig() {
	configFileName := "config.yml"
	_, err := os.Stat(configFileName)
	if err == nil {
		slog.Info("InitChooseConfig", "配置文件选择", "本地配置文件")
		InitConfig()
	}
	if os.IsNotExist(err) {
		slog.Info("InitChooseConfig", "配置文件选择", "环境变量")
		InitEnv()
	}
	slog.Info("InitChooseConfig", "config", global.ServerConfig)
}

func InitConfig() {
	configFileName := "config.yml"
	v := viper.New()
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		slog.Error("initialize", "ReadInConfig error", err.Error())
		os.Exit(1)
	}
	if err := v.Unmarshal(global.ServerConfig); err != nil {
		slog.Error("initialize", "Unmarshal error", err.Error())
		os.Exit(1)
	}
	slog.Info("initialize", "Config", global.ServerConfig)
	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		// 加锁，确保更新配置的过程中其他地方无法访问
		configLock.Lock()
		defer configLock.Unlock()
		newConfig := &config.ServerConfig{}
		if err := v.ReadInConfig(); err != nil {
			slog.Error("config file changed error", err)
		}
		if err := v.Unmarshal(newConfig); err != nil {
			slog.Error("config file changed  error", err)
		}
		global.ServerConfig = newConfig
	})
}

func InitDiscord() {
	discord, err := gpt.InitDiscord()
	if err != nil {
		slog.Error("initialize", "InitDiscord error", err.Error())
	}
	global.DiscordSession = discord
}

func InitEnv() {
	viper.AutomaticEnv()
	cnf := config.ServerConfig{}
	channelId := viper.GetString("CHANNEL_ID")
	botId := viper.GetString("BOT_ID")
	sendBotToken := viper.GetString("SEND_BOT_TOKEN")
	proxyUrl := viper.GetString("PROXY_URL")
	proxyEnable := viper.GetBool("PROXY_ENABLE")
	discordConfig := config.DiscordConfig{
		ChannelId:    channelId,
		BotId:        botId,
		SendBotToken: sendBotToken,
		ProxyEnable:  proxyEnable,
		ProxyUrl:     proxyUrl,
	}
	cnf.DiscordConfig = discordConfig
	pushUrl := viper.GetString("PUSH_URL")
	pushToken := viper.GetString("PUSH_TOKEN")
	pushPlus := config.PushPlus{
		Token: pushToken,
		Url:   pushUrl,
	}
	cnf.PushPlus = pushPlus
	apiBase := viper.GetString("ONE_API_BASE")
	apiToken := viper.GetString("API_TOKEN")
	apiConfig := config.OneApiConfig{
		Proxy:  apiBase,
		SToken: apiToken,
	}
	cnf.OneApiConfig = apiConfig
	imagePrefix := viper.GetString("MODE_IMAGE_PREFIX")
	mode := config.Mode{ImagePrefix: imagePrefix}
	cnf.Mode = mode
	chatProxyEnable := viper.GetBool("CHAT_PROXY_ENABLE")
	chatProxyUrl := viper.GetString("CHAT_PROXY_URL")
	autoPass := viper.GetBool("AUTO_PASS")
	chat := config.Chat{
		Proxy:          chatProxyEnable,
		ProxyUrl:       chatProxyUrl,
		SessionTimeOut: 60,
		Model:          "gpt-4",
		AutoPass:       autoPass,
	}
	cnf.Chat = chat
	global.ServerConfig = &cnf
}
