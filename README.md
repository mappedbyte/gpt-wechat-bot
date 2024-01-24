<div align="center">
<h1>GPT Wechat Bot </h1>
<p>  🎨基于GO语言实现的微信聊天机器人🎨 </p>
</div><div align="left"></div>
个人微信接入ChatGPT，实现和GPT机器人互动聊天，支持私聊回复和群聊艾特回复。


### 实现功能

* GPT机器人模型可配置
* 支持gpt3.5,gpt4模型
* 私聊支持上下文
* 机器人私聊回复&机器人群聊@回复
* 好友添加自动通过可配置
* 机器人掉线触发微信公众号推送

### 暂未实现
* 图片生成
* 个性化指令定制


### 实现机制
1. 利用微信A作为机器人扫码登录程序模拟的微信电脑端，程序后端调用API接口进行文本回复和图片生成。其他微信账号与微信A聊天实现微信个人机器人功能。基于[openwechat](https://github.com/eatmoreapple/openwechat)开源仓库实现

> GPT的[官方文档](https://beta.openai.com/docs/models/overview)和详细[参数示例](https://beta.openai.com/examples) 。
>


### 注意事项

* 项目仅供娱乐，滥用可能有微信封禁的风险，请勿用于商业用途。
* 未对敏感词汇进行过滤，如需过滤请自行添加


### 结果展示

#### 个人聊天
<img src="image/use_msg.png"/>

#### 群聊@回复
<img src="image/group_msg.png"/>



### 使用说明

```
chat:
  autoPass: true #是否自动通过好友
  proxy: false #是否使用代理
  proxyUrl: http://127.0.0.1:7890 #代理地址
  sessionTimeOut: 60
  model: gpt-4  #替换模型

one-api:
  proxy: #替换为openAI或第三方的API地址
  s-token: #替换为接口的token

push:
  url: http://www.pushplus.plus/send
  token: #替换为自己的token

```