package main

import (
	chat_pipeline2 "DocReader/internal/agent/chat_pipeline"
	"DocReader/pkg/mem"
	"context"
	"fmt"

	"github.com/cloudwego/eino/schema"
)

func main() {
	ctx := context.Background()
	id := "test-user-001"

	// Build chat agent
	runner, err := chat_pipeline2.BuildChatAgent(ctx)
	if err != nil {
		panic(err)
	}

	// First conversation
	fmt.Println("========== First Conversation ==========")
	userMessage := &chat_pipeline2.UserMessage{
		ID:      id,
		Query:   "Hello, what can you help me with?",
		History: mem.GetSimpleMemory(id).GetMessages(),
	}

	out, err := runner.Invoke(ctx, userMessage)
	if err != nil {
		panic(err)
	}
	fmt.Println("Q: Hello, what can you help me with?")
	fmt.Println("A:", out.Content)

	// Save to memory
	mem.GetSimpleMemory(id).SetMessages(schema.UserMessage("Hello, what can you help me with?"))
	mem.GetSimpleMemory(id).SetMessages(schema.SystemMessage(out.Content))

	// Second conversation (test context memory)
	fmt.Println("\n========== Second Conversation ==========")
	userMessage = &chat_pipeline2.UserMessage{
		ID:      id,
		Query:   "What is the current time?",
		History: mem.GetSimpleMemory(id).GetMessages(),
	}

	out, err = runner.Invoke(ctx, userMessage)
	if err != nil {
		panic(err)
	}
	fmt.Println("Q: What is the current time?")
	fmt.Println("A:", out.Content)

	fmt.Println("\n========== Chat Test Completed ==========")
}
