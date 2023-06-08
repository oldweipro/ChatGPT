package chat

import (
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"net/http"
	"net/url"
)

func EmbeddingFunc() {
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
	client := openai.NewClientWithConfig(config)
	ctx := context.Background()
	hello := []string{
		"go",
		"时",
		"间",
		"格",
		"式",
		"化",
	}
	//hello := []string{
	//	"你",
	//	"好",
	//}
	req := openai.EmbeddingRequest{
		Model: openai.AdaEmbeddingV2,
		Input: hello,
	}
	stream, err := client.CreateEmbeddings(ctx, req)
	if err != nil {
		fmt.Printf("CreateEmbeddings error: %v\n", err)
		return
	}
	fmt.Println(len(stream.Data))
}
