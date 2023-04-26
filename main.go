package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"io"
	"net/http"
	"net/url"
	"os"
)

var messages = []openai.ChatCompletionMessage{
	{
		Role:    openai.ChatMessageRoleSystem,
		Content: "你是由openfish训练的人工智能助手",
	},
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
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
			ChatCompletionStream()
		} else {
			// 如果出现了错误，则退出循环
			break
		}
	}
	// 输出结束信息
	fmt.Println("程序已结束")
}

func ChatCompletionStream() {
	config := openai.DefaultConfig("YOUR_OPENAI_KEY")
	// 如果需要代理，请配置代理地址，如不需要可注释或删掉以下代码
	config.HTTPClient.Transport = &http.Transport{
		// 设置Transport字段为自定义Transport，包含代理设置
		Proxy: func(req *http.Request) (*url.URL, error) {
			// 设置代理
			proxyURL, err := url.Parse("http://127.0.0.1:7890")
			if err != nil {
				return nil, err
			}
			return proxyURL, nil
		},
	}
	client := openai.NewClientWithConfig(config)
	ctx := context.Background()

	req := openai.ChatCompletionRequest{
		Model:     openai.GPT3Dot5Turbo0301,
		MaxTokens: 1000,
		Messages:  messages,
		Stream:    true,
	}
	stream, err := client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		fmt.Printf("ChatCompletionStream error: %v\n", err)
		return
	}
	defer stream.Close()

	fmt.Printf("回答: ")
	var streamResponse string
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println()
			messages = append(messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: streamResponse,
			})
			return
		}
		if err != nil {
			return
		}
		streamResponse += response.Choices[0].Delta.Content
		fmt.Printf(response.Choices[0].Delta.Content)
	}
}
