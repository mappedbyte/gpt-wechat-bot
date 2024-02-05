package gpt

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"log/slog"
	"net/http"
	"net/url"
	"time"
	"wechatBot/internal/global"
)

var UserMention string
var ResponseImage = make(map[string][]string)

type ImageMessage struct {
	MessageId string
	Image     []string
}

func InitDiscord() (*discordgo.Session, error) {
	enable := global.ServerConfig.DiscordConfig.ProxyEnable
	discord, err := discordgo.New("Bot " + global.ServerConfig.DiscordConfig.SendBotToken)
	if enable {
		if err != nil {
			slog.Error("InitDiscord", "errorMsg", err.Error())
			return nil, err
		}
		proxyURL, err := url.Parse(global.ServerConfig.DiscordConfig.ProxyUrl)
		if err != nil {
			slog.Error("InitDiscord", "errorMsg", "错误的代理地址:"+err.Error())
			return nil, err
		}
		discord.Dialer.Proxy = http.ProxyURL(proxyURL)
		discord.Client = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)}}
	}
	discord.AddHandler(messageCreate)
	err = discord.Open()
	if err != nil {
		slog.Error("InitDiscord", "errorMsg", err.Error())
		return nil, err
	}
	return discord, nil
}

func ReplyImage(imagePrompt string) []string {
	/*	if !strings.HasPrefix(imagePrompt, "画") {
		slog.Info("ReplyImage", "imagePrompt", imagePrompt)
		imagePrompt = fmt.Sprintf(` {"model": "dall-e-3","prompt": "%s","n": 1,"size": "1024x1024"}`, imagePrompt)
	}*/
	imagePrompt = " " + imagePrompt
	images := make([]string, 0)
	//images = append(images, global.DeadlineExceededImage)
	if UserMention == "" {
		user, err := global.DiscordSession.User(global.ServerConfig.DiscordConfig.BotId)
		if err != nil {
			slog.Error("ReplyImage", "获取bot用户失败", err.Error())
			images = append(images, global.DeadlineExceededImage)
			return images
		}
		UserMention = user.Mention()
	}
	message, err := global.DiscordSession.ChannelMessageSend(global.ServerConfig.DiscordConfig.ChannelId, UserMention+imagePrompt)
	if err != nil {
		images = append(images, global.DeadlineExceededImage)
		return images
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), 600*time.Second)
	slog.Info("ReplyImage", "请求的接口MessageId", message.ID)
	go WatchImage(ctx, cancelFunc, message.ID)
	<-ctx.Done()
	return ResponseImage[message.ID]
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	go func() {
		slog.Info("messageCreate", "返回的messageId", m.ID)
		if m.ReferencedMessage != nil {
			images := make([]string, 0)
			slog.Info("messageCreate", "返回的引用的messageId", m.ReferencedMessage.ID)
			for _, embed := range m.Embeds {
				if embed != nil {
					pic := embed.Image.URL
					images = append(images, pic)
				}
			}

			if len(images) == 0 {
				ticker := time.NewTicker(1 * time.Second)
				ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
				go Operation(ctx, cancel, ticker, s, m)
				<-ctx.Done()
				slog.Info("messageCreate", "go Operation Done", "运行结束")
				message, err := s.ChannelMessage(m.ChannelID, m.Message.ID)
				if err != nil {
					slog.Error("messageCreate", "重新获取图片状态失败", err.Error())
				}

				slog.Info("messageCreate", "message.Embeds", message.Embeds)

				for _, embed := range message.Embeds {
					if embed.Image != nil {
						slog.Info("messageCreate", "embed.Image", embed.Image)
						images = append(images, embed.Image.URL)
					}
				}
				if len(images) == 0 {
					images = append(images, global.DeadlineExceededImage)
				}
				slog.Info("ChannelMessage", "返回的图片", images)
			}
			ResponseImage[m.ReferencedMessage.ID] = images
			slog.Info("Message", "返回的图片", images)
		}
	}()

	/*	//只集成了图片信息,所以要等待图片,不考虑聊天信息
		if len(m.ReferencedMessage.Embeds) == 0 {
			ticker := time.NewTicker(1 * time.Second)
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			go Operation(ctx, cancel, ticker, s, m)
			<-ctx.Done()
			message, err := s.ChannelMessage(m.ChannelID, m.ID)
			if err != nil {
				slog.Error("messageCreate", "重新获取图片状态失败", err.Error())
			}
			for _, embed := range message.Embeds {
				if embed != nil {
					images = append(images, embed.URL)
				}
			}
			slog.Info("ChannelMessage", "返回的图片", images)
		}

		ResponseImage[m.ID] = images
	}*/

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Pong!")
	}
	// If the message is "pong" reply with "Ping!"
	if m.Content == "pong" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Ping!")
	}

}

func Operation(ctx context.Context, cancel context.CancelFunc, ticker *time.Ticker, s *discordgo.Session, m *discordgo.MessageCreate) {
	defer cancel()
	for {
		select {
		case <-ticker.C:
			message, err := s.ChannelMessage(m.ChannelID, m.Message.ID)
			if err == nil {
				if message.Embeds != nil && len(message.Embeds) > 0 {
					ticker.Stop()
					return
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

func WatchImage(ctx context.Context, cancel context.CancelFunc, messageId string) {
	defer cancel()
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if ResponseImage[messageId] != nil && len(ResponseImage[messageId]) > 0 {
				ticker.Stop()
				return
			}
		}
	}
}
