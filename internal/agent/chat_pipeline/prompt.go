package chat_pipeline

import (
	"context"

	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

type ChatTemplateConfig struct {
	FormatType schema.FormatType
	Templates  []schema.MessagesTemplate
}

// newChatTemplate component initialization function of node 'ChatTemplate' in graph 'ChatAgent'
func newChatTemplate(ctx context.Context) (ctp prompt.ChatTemplate, err error) {
	config := &ChatTemplateConfig{
		FormatType: schema.FString,
		Templates: []schema.MessagesTemplate{
			schema.SystemMessage(systemPrompt),
			schema.MessagesPlaceholder("history", false),
			schema.UserMessage("{content}"),
		},
	}
	ctp = prompt.FromMessages(config.FormatType, config.Templates...)
	return ctp, nil
}

var systemPrompt = `
# Role: Intelligent Document Assistant

## Core Capabilities
- Answer questions based on knowledge base documents
- Context understanding and multi-turn conversation
- Provide accurate and helpful responses

## Guidelines
- Prioritize answering based on document content
- If the information is not found in documents, honestly inform the user
- Provide clear, concise, and well-structured responses
- When appropriate, cite relevant parts of the documents

## Response Requirements
- Be accurate and factual
- Structure responses for readability
- Use plain text format (avoid markdown syntax)
- If uncertain, acknowledge limitations

## Context Information
- Current date: {date}
- Related documents:
==== Documents Start ====
{documents}
==== Documents End ====
`
