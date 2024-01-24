package main

import (
	"wechatBot/internal/application"
)

func main() {
	bot := application.Run()
	_ = bot.Block()
}
