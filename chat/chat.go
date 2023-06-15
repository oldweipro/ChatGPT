package chat

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
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
var text = "你好讯飞"

const Base64 = ""

var functions = `[{
"name": "get_current_weather",
"description": "Get the current weather",
"parameters": {
"type": "object",
"properties": {
"location": {
"type": "string",
"description": "The city and state, e.g. San Francisco, CA",
},
"format": {
"type": "string",
"enum": ["celsius", "fahrenheit"],
"description": "The temperature unit to use. Infer this from the users location.",
},
},
"required": ["location", "format"],
},
}]`

func RunCmdChatGPT(f func()) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("%s 提问: ", time.Now().Format("2006-01-02 15:04:05"))
		if scanner.Scan() {
			// 获取用户输入的文本
			text = scanner.Text()
			// 打印用户输入的文本sssjjjjdhjksahk
			fmt.Println("")
			marshal, _ := json.Marshal(functions)
			messages = append(messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: text,
			})
			messages = append(messages, openai.ChatCompletionMessage{
				Role: openai.ChatMessageRoleAssistant,
				Name: "get_current_weather",
				FunctionCall: &openai.FunctionCall{
					Name:      "get_current_weather",
					Arguments: string(marshal),
				},
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

func CompletionStreamKey() {
	config := openai.DefaultConfig("")
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
		Model:     openai.GPT3Dot5Turbo16K0613,
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

const Cookies = ""
const GtToken = "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"

func XinghuoChatMessage() {
	client := resty.New()
	formData := make(map[string]string)
	formData["chatId"] = "5277689"
	formData["GtToken"] = GtToken
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
			Value: Cookies,
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
	formData["chatId"] = "5320200"
	formData["GtToken"] = GtToken
	// 正常对话不需要填写。如果使用助手值为1
	formData["isBot"] = "0"
	// 助手ID
	formData["botId"] = "0"
	// app
	formData["clientType"] = "4"
	formData["text"] = text
	//vcn params is empty;code=11119是什么错误
	resp, _ := client.R().
		//SetHeader("Accept", "text/event-stream").
		//SetHeader("User-Agent", "okhttp/4.11.0").
		//SetHeader("Host", "xinghuo.xfyun.cn").
		//SetHeader("app", "xinghuo").
		//SetHeader("clientType", "4").
		//SetHeader("platform", "android").
		//SetHeader("versionCode", "2023061002").
		//SetHeader("Accept-Encoding", "gzip").
		SetCookie(&http.Cookie{
			Name:  "ssoSessionId",
			Value: Cookies,
		}).
		SetFormData(formData).SetDoNotParseResponse(true).
		Post("https://xinghuo.xfyun.cn/iflygpt-app/u/chat_message/chat")
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
