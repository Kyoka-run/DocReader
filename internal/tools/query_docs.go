package tools

import (
	"DocReader/internal/component/retriever"
	"context"
	"encoding/json"
	"log"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
)

type QueryDocsInput struct {
	Query string `json:"query" jsonschema:"description=The query string to search in knowledge base for relevant documents"`
}

func NewQueryDocsTool() tool.InvokableTool {
	t, err := utils.InferOptionableTool(
		"query_docs",
		"Search the knowledge base for relevant documents. Use this tool when you need to find information from uploaded documents to answer user questions.",
		func(ctx context.Context, input *QueryDocsInput, opts ...tool.Option) (output string, err error) {
			rr, err := retriever.NewMilvusRetriever(ctx)
			if err != nil {
				log.Printf("Failed to create retriever: %v", err)
				return "", err
			}
			resp, err := rr.Retrieve(ctx, input.Query)
			if err != nil {
				log.Printf("Failed to retrieve documents: %v", err)
				return "", err
			}
			respBytes, _ := json.Marshal(resp)
			output = string(respBytes)
			return output, nil
		})
	if err != nil {
		log.Fatal(err)
	}
	return t
}
