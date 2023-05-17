package main

import (
	"openai/gin_server"
)

func main() {
	//chat.RunCmdChatGPT(chat.CompletionStreamReverse)
	gin_server.GinServer()
}
