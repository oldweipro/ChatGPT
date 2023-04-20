package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"io"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "你是由openfish训练的人工智能助手",
		},
	}
	for {
		fmt.Print("提问：")
		if scanner.Scan() {
			// 获取用户输入的文本
			text := scanner.Text()
			// 打印用户输入的文本
			messages = append(messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: text,
			})
			ais := run(messages)
			messages = append(messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: ais,
			})
		} else {
			// 如果出现了错误，则退出循环
			break
		}
	}
	// 输出结束信息
	fmt.Println("程序已结束")
}
func run(messages []openai.ChatCompletionMessage) string {
	c := openai.NewClient("YOUR_OPENAI_KEY")
	ctx := context.Background()

	req := openai.ChatCompletionRequest{
		Model:     openai.GPT3Dot5Turbo,
		MaxTokens: 1000,
		Messages:  messages,
		Stream:    true,
	}
	stream, err := c.CreateChatCompletionStream(ctx, req)
	if err != nil {
		fmt.Printf("ChatCompletionStream error: %v\n", err)
		return ""
	}
	defer stream.Close()

	fmt.Printf("回答: ")
	streamResponse := ""
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println()
			return streamResponse
		}

		if err != nil {
			return ""
		}
		streamResponse += response.Choices[0].Delta.Content
		fmt.Printf(response.Choices[0].Delta.Content)
	}
}
