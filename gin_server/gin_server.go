package gin_server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	"net/http"
	"openai/chat"
	"sync"
	"time"
)

var messages = []openai.ChatCompletionMessage{
	{
		Role:    openai.ChatMessageRoleSystem,
		Content: "你是由openfish训练的人工智能助手",
	},
}

func GinServer() {
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	})
	r.POST("/chat-process", ChatProcess)
	r.GET("/api", HandleAPI)
	err := r.Run(":8787")
	if err != nil {
		return
	}
}

var (
	mutex         sync.Mutex // 并发控制锁
	counter       int        // 计数器
	lastTimestamp time.Time  // 上一个时间戳
)

func OneMinuteThree(c *gin.Context) {
	// 获取互斥锁
	mutex.Lock()
	defer mutex.Unlock()

	// 检查是否需要重置计数器
	if time.Since(lastTimestamp) > time.Minute {
		counter = 0
	}

	// 检查并发是否达到上限
	if counter >= 2 {
		c.JSON(429, gin.H{
			"error": "请求过多",
		})
		return
	}

	// 更新计数器和时间戳
	counter++
	lastTimestamp = time.Now()

	// 在这里执行需要同步的代码
	// ...

	// 返回结果
	c.JSON(200, gin.H{
		"message": "请求已处理",
	})
}

var semaphore = make(chan struct{}, 3)

func HandleAPI(c *gin.Context) {
	semaphore <- struct{}{}
	defer func() { <-semaphore }()

	c.JSON(200, gin.H{
		"message": "请求已处理",
	})
}

func ChatProcess(c *gin.Context) {
	var chatMsg openai.ChatCompletionMessage
	err := c.ShouldBindJSON(&chatMsg)
	if err != nil {
		fmt.Println(err)
		return
	}
	c.Status(http.StatusOK)
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	messages = append(messages, chatMsg)
	chat.CompletionStreamReverse()
}
