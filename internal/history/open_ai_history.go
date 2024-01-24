package history

import (
	"github.com/eatmoreapple/openwechat"
	"github.com/patrickmn/go-cache"
	"time"
	"wechatBot/internal/gpt"
)

type OpenAIHistory struct {
	user  *openwechat.User
	cache *cache.Cache
}

func (u *OpenAIHistory) ClearUserHistory() {
	u.cache.Delete(u.user.AvatarID())
	history := make([]any, 0)
	//这个可以根据自己的AI平台去设置一个默认的
	history = append(history, gpt.Message{Role: "system",
		Content: "\nYou are ChatGPT, a large language model trained by OpenAI.\nKnowledge cutoff: 2021-09\nCurrent model: gpt-4\nCurrent time: 2024/1/12 19:18:07\nLatex inline: $x^2$ \nLatex block: $$e=mc^2$$\n\n",
	})
	u.cache.Set(u.user.AvatarID(), history, 60*time.Minute)
}

func (u *OpenAIHistory) GetHistory() []any {
	historyList, ok := u.cache.Get(u.user.AvatarID())
	var anyList []interface{}
	if !ok {
		history := make([]any, 0)
		//这个可以根据自己的AI平台去设置一个默认的
		history = append(history, gpt.Message{Role: "system",
			Content: "\nYou are ChatGPT, a large language model trained by OpenAI.\nKnowledge cutoff: 2021-09\nCurrent model: gpt-4\nCurrent time: 2024/1/12 19:18:07\nLatex inline: $x^2$ \nLatex block: $$e=mc^2$$\n\n",
		})
		u.cache.Set(u.user.AvatarID(), history, 60*time.Minute)
		for _, s := range history {
			anyList = append(anyList, s)
		}
		return anyList
	}
	historyArray := historyList.([]any)
	if len(historyArray) >= 100 {
		defer u.ClearUserHistory()
	}
	for _, s := range historyArray {
		anyList = append(anyList, s)
	}
	return anyList
}

func (u *OpenAIHistory) SetUserHistory(text any) {

	historyList, ok := u.cache.Get(u.user.AvatarID())
	var anyList []interface{}
	if !ok {
		anyList = append(anyList, text)
		u.cache.Set(u.user.AvatarID(), anyList, 60*time.Minute)
		return
	}
	historyArray := historyList.([]any)
	historyArray = append(historyArray, text)
	u.cache.Set(u.user.AvatarID(), historyArray, 60*time.Minute)
}
func NewOpenAIUserHistory(user *openwechat.User, cache *cache.Cache) *OpenAIHistory {
	return &OpenAIHistory{
		user:  user,
		cache: cache,
	}
}
