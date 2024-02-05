package main

import (
	"context"
	"log/slog"
	"time"
	"wechatBot/internal/application"
)

func main() {
	bot := application.Run()
	err := bot.Block()
	slog.Error("main", "bot err", err.Error())
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	<-ctx.Done()
	slog.Info("main", "exit", "程序退出!")
}
