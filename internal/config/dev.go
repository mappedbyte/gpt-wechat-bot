package config

import (
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/skip2/go-qrcode"
	"log"
	"runtime"
)

func CheckOs() func(uuid string) {
	if runtime.GOOS == "darwin" {
		log.Println("当前是 Mac 系统")
		return openwechat.PrintlnQrcodeUrl
	} else if runtime.GOOS == "linux" {
		log.Println("当前是 Linux 系统")
		return ConsoleQrCode
	} else if runtime.GOOS == "windows" {
		log.Println("当前是 Windows 系统")
		return openwechat.PrintlnQrcodeUrl
	}
	log.Fatalln("无法确定当前系统")
	return nil
}

func ConsoleQrCode(uuid string) {
	fmt.Println("扫描控制台二维码登录")
	q, _ := qrcode.New("https://login.weixin.qq.com/l/"+uuid, qrcode.Low)
	fmt.Println(q.ToString(true))
}
