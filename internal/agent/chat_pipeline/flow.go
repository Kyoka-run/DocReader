package chat_pipeline

import (
	tools2 "DocReader/internal/tools"
	"context"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
)

func newReactAgentLambda(ctx context.Context) (lba *compose.Lambda, err error) {
	config := &react.AgentConfig{
		MaxStep:            25,
		ToolReturnDirectly: map[string]struct{}{},
	}

	chatModelIns, err := newChatModel(ctx)
	if err != nil {
		return nil, err
	}
	config.ToolCallingModel = chatModelIns

	// Only keep document-related tools
	config.ToolsConfig.Tools = append(config.ToolsConfig.Tools, tools2.NewGetCurrentTimeTool())
	config.ToolsConfig.Tools = append(config.ToolsConfig.Tools, tools2.NewQueryDocsTool())

	ins, err := react.NewAgent(ctx, config)
	if err != nil {
		return nil, err
	}

	lba, err = compose.AnyLambda(ins.Generate, ins.Stream, nil, nil)
	if err != nil {
		return nil, err
	}
	return lba, nil
}
