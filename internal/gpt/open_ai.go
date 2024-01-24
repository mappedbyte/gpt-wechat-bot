package gpt

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"wechatBot/internal/global"
)

type Response struct {
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Delta Delta `json:"delta"`
}

type Delta struct {
	Content string `json:"content"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type OpenAI struct {
}

func (o *OpenAI) Chat(msg []any) (string, error) {
	params := map[string]interface{}{
		"messages":          msg,
		"stream":            true,
		"model":             global.ServerConfig.Chat.Model,
		"temperature":       0.5,
		"presence_penalty":  0,
		"frequency_penalty": 0,
		"top_p":             1,
	}
	body, _ := json.Marshal(params)
	fmt.Println(string(body))
	request, _ := http.NewRequest("POST", global.ServerConfig.OneApiConfig.Proxy+"/v1/chat/completions", bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+global.ServerConfig.OneApiConfig.SToken)
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	response, err := global.Client.Do(request)
	fmt.Println(response.StatusCode)
	if err != nil || response.StatusCode != 200 {
		slog.Info("OpenAI", "请求gpt接口出现异常", "responseStatus:"+strconv.Itoa(response.StatusCode))
		return global.DeadlineExceededText, errors.New("请求gpt接口出现异常")
	}
	slog.Info("OpenAI", "GPT Response Status", strconv.Itoa(response.StatusCode))
	defer response.Body.Close()
	var responseText string
	reader := bufio.NewReader(response.Body)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Printf("Read error: %s\n", err)
			break
		}
		lineString := fmt.Sprintf("Received: %s", line)
		if len(lineString) > 0 {
			jsonData := strings.SplitN(lineString, "Received: data: ", 2)
			if len(jsonData) >= 2 {
				var res Response
				fmt.Println(jsonData, len(jsonData))
				err := json.Unmarshal([]byte(jsonData[1]), &res)
				if err != nil {
					continue
				}
				for _, choice := range res.Choices {
					responseText += choice.Delta.Content
				}
			}
		}
	}
	return responseText, nil
}
