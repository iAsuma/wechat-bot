package service

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/sashabaranov/go-openai"
	"wechatbot/internal/model"
)

type lOpenAiService struct{}

func NewOpenAiService() *lOpenAiService {
	return &lOpenAiService{}
}

func (l *lOpenAiService) Chat(ctx context.Context, inputText string) string {

	config, err := g.Cfg().Get(ctx, "openai")
	if config.IsNil() || err != nil {
		return ""
	}
	openaiConfig := model.OpenConfig{}
	_ = config.Struct(&openaiConfig)

	reqConfig := openai.DefaultConfig(openaiConfig.ApiKey)
	reqConfig.BaseURL = openaiConfig.BaseUrl

	client := openai.NewClientWithConfig(reqConfig)
	reqModel := openai.GPT3Dot5Turbo
	if openaiConfig.Model != "" {
		reqModel = openaiConfig.Model
	}

	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: reqModel,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "你是万能小助手",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: inputText,
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return ""
	}

	return resp.Choices[0].Message.Content
}
