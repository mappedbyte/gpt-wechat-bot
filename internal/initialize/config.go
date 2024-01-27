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
