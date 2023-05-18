package main

import (
	"github.com/oldweipro/chatgpt/gin_server"
)

func main() {
	//chat.RunCmdChatGPT(chat.CompletionStreamReverse)
	gin_server.GinServer()
}
