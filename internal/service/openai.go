package service

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/sashabaranov/go-openai"
	"wechatbot/internal/model"
)

type lOpenAiService struct{}

func NewOpenAi() *lOpenAiService {
	return &lOpenAiService{}
}

func (l *lOpenAiService) Chat(ctx context.Context, msgList []model.ChatMessage) string {
	config, err := g.Cfg().GetWithEnv(ctx, "openai")
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

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "你是万能小助手",
		},
	}

	for _, v := range msgList {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    v.Role,
			Content: v.Content,
		})
	}

	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:    reqModel,
			Messages: messages,
		},
	)

	if err != nil {
		errMsg := fmt.Sprintf("ChatCompletion error: %v\n", err)
		g.Log().Error(ctx, errMsg)
		NewEmail().Send(context.Background(), "WX-Bot Gpt Error", errMsg)
		return ""
	}

	return resp.Choices[0].Message.Content
}

func (l *lOpenAiService) Ask(ctx context.Context, inputText string) string {
	if inputText == "" {
		return ""
	}

	messages := make([]model.ChatMessage, 0)
	messages = append(messages, model.ChatMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: inputText,
	})

	return l.Chat(ctx, messages)
}
