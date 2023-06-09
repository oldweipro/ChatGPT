package chat

import (
	"bufio"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/sashabaranov/go-openai"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var messages []openai.ChatCompletionMessage
var text string

func RunCmdChatGPT(f func()) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("%s 提问: ", time.Now().Format("2006-01-02 15:04:05"))
		if scanner.Scan() {
			// 获取用户输入的文本
			text = scanner.Text()
			// 打印用户输入的文本
			messages = append(messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: text,
			})
			fmt.Printf("%s 回答: ", time.Now().Format("2006-01-02 15:04:05"))
			f()
		} else {
			// 如果出现了错误，则退出循环
			break
		}
	}
	// 输出结束信息
	fmt.Println("程序已结束")
}

func CompletionStreamKey(openaiKey string) {
	config := openai.DefaultConfig(openaiKey)
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
	Base(config)
}

func CompletionStreamReverse() {
	config := openai.DefaultConfig("")
	config.BaseURL = "http://127.0.0.1:8080/v1"
	Base(config)
}

func Base(config openai.ClientConfig) {
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

func XinghuoChatMessage() {
	client := resty.New()
	formData := make(map[string]string)
	formData["chatId"] = "4422137"
	formData["GtToken"] = ""
	formData["fd"] = ""
	formData["clientType"] = "2"
	formData["text"] = text
	//vcn params is empty;code=11119是什么错误
	resp, _ := client.R().
		SetHeader("Accept", "text/event-stream").
		SetHeader("Referer", "https://xinghuo.xfyun.cn/chat?id="+formData["chatId"]).
		SetHeader("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 16_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.5 Mobile/15E148 Safari/604.1").
		SetCookie(&http.Cookie{
			Name:  "ssoSessionId",
			Value: "",
		}).
		SetFormData(formData).SetDoNotParseResponse(true).
		Post("https://xinghuo.xfyun.cn/iflygpt-chat/u/chat_message/chat")
	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.RawBody())

	buf := make([]byte, 1024)
	textContent := ""
	for {
		n, err := resp.RawBody().Read(buf)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
		if n == 0 {
			break
		}
		msg := string(buf[:n])
		if strings.Contains(msg, "<end>") {
			break
		}
		result := processMsg(msg)
		decodeString, msgErr := base64.StdEncoding.DecodeString(result)
		if msgErr != nil {
			fmt.Println(msg)
			fmt.Println("base64解码失败")
		}
		textContent += string(decodeString)
		fmt.Print(string(decodeString))
	}
	fmt.Println()
}
func XinghuoChatMessageApp() {
	client := resty.New()
	formData := make(map[string]string)
	formData["chatId"] = "4582237"
	formData["GtToken"] = ""
	// 正常对话不需要填写。如果使用助手值为1
	formData["isBot"] = "1"
	// 助手ID
	formData["botId"] = "1626"
	// app
	formData["clientType"] = "4"
	formData["text"] = text
	//vcn params is empty;code=11119是什么错误
	resp, _ := client.R().
		SetHeader("Accept", "text/event-stream").
		SetHeader("User-Agent", "okhttp/4.11.0").
		SetHeader("Host", "xinghuo.xfyun.cn").
		SetHeader("app", "xinghuo").
		SetHeader("clientType", "4").
		SetHeader("platform", "android").
		SetHeader("versionCode", "2023061002").
		SetHeader("Accept-Encoding", "gzip").
		SetCookie(&http.Cookie{
			Name:  "ssoSessionId",
			Value: "",
		}).
		SetFormData(formData).SetDoNotParseResponse(true).
		Post("https://xinghuo.xfyun.cn/iflygpt-chat/u/chat_message/chat")
	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.RawBody())

	buf := make([]byte, 1024)
	textContent := ""
	for {
		n, err := resp.RawBody().Read(buf)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
		if n == 0 {
			break
		}
		msg := string(buf[:n])
		if strings.Contains(msg, "<end>") {
			break
		}
		result := processMsg(msg)
		decodeString, msgErr := base64.StdEncoding.DecodeString(result)
		if msgErr != nil {
			fmt.Println(msg)
			fmt.Println("base64解码失败")
		}
		textContent += string(decodeString)
		fmt.Print(string(decodeString))
	}
	fmt.Println()
}

func processMsg(msg string) string {
	if strings.Contains(msg, "data:") {
		for _, str := range strings.Split(msg, "\n\n") {
			if len(str) > 0 {
				if strings.HasPrefix(str, "data:") {
					return str[5:]
				} else {
					return str
				}
			}
		}
	} else {
		return msg
	}
	return ""
}

func Loop() {
	// Create a goroutine to run the task every 10 minutes
	go func() {
		for {
			// Print "hello world"
			text = "行政人事"
			fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
			XinghuoChatMessageApp()
			// Sleep for 10 minutes
			time.Sleep(10 * time.Minute)
		}
	}()

	// Keep the main goroutine running indefinitely
	select {}
}
