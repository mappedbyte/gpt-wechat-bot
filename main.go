package main

import (
	"wechatBot/internal/application"
)

func main() {
	bot := application.Run()
	application.RunGin()
	_ = bot.Block()
}
