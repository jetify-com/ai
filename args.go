package aisdk

import (
	"fmt"

	"go.jetify.com/ai/api"
	"go.jetify.com/pkg/try"
)

type llmArgs struct {
	Prompt []api.Message
	Config GenerateTextConfig
}

func toLLMArgs(args ...any) (*llmArgs, error) {
	result := &llmArgs{
		Prompt: []api.Message{},
		Config: GenerateTextConfig{},
	}

	var pendingBlocks []api.ContentBlock

	for _, arg := range args {
		if err := processArg(arg, result, &pendingBlocks); err != nil {
			return nil, err
		}
	}

	// Make sure to flush any remaining blocks
	if err := flushBlocks(result, &pendingBlocks); err != nil {
		return nil, err
	}

	return result, nil
}

func processArg(arg any, result *llmArgs, blocks *[]api.ContentBlock) error {
	switch typedArg := arg.(type) {
	case string:
		*blocks = append(*blocks, &api.TextBlock{Text: typedArg})
	case api.ContentBlock:
		*blocks = append(*blocks, typedArg)
	case *api.ContentBlock:
		*blocks = append(*blocks, *typedArg)
	case []api.ContentBlock:
		*blocks = append(*blocks, typedArg...)
	case try.Try[api.ContentBlock]:
		block, err := typedArg.Get()
		if err != nil {
			return err
		}
		*blocks = append(*blocks, block)
	case try.Try[*api.ContentBlock]:
		block, err := typedArg.Get()
		if err != nil {
			return err
		}
		*blocks = append(*blocks, *block)
	case api.Message, []api.Message:
		if err := flushBlocks(result, blocks); err != nil {
			return err
		}
		addToPrompt(result, typedArg)
	case try.Try[api.Message]:
		msg, err := typedArg.Get()
		if err != nil {
			return err
		}
		if err := flushBlocks(result, blocks); err != nil {
			return err
		}
		addToPrompt(result, msg)
	case GenerateOption:
		typedArg(&result.Config)
	case []GenerateOption:
		for _, opt := range typedArg {
			opt(&result.Config)
		}
	default:
		return fmt.Errorf("unsupported argument type: %T", arg)
	}
	return nil
}

func flushBlocks(result *llmArgs, blocks *[]api.ContentBlock) error {
	if len(*blocks) > 0 {
		result.Prompt = append(result.Prompt, &api.UserMessage{
			Content: *blocks,
		})
		*blocks = nil
	}
	return nil
}

func addToPrompt(result *llmArgs, arg any) {
	switch v := arg.(type) {
	case api.Message:
		result.Prompt = append(result.Prompt, v)
	case []api.Message:
		result.Prompt = append(result.Prompt, v...)
	}
}

// processContentArgs processes arguments to build content blocks
func processContentArgs(args ...any) ([]api.ContentBlock, *api.ProviderMetadata, error) {
	var blocks []api.ContentBlock
	var metadata *api.ProviderMetadata

	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			blocks = append(blocks, &api.TextBlock{Text: v})
		case api.ContentBlock:
			blocks = append(blocks, v)
		case []api.ContentBlock:
			blocks = append(blocks, v...)
		case *api.ProviderMetadata:
			if metadata != nil {
				return nil, nil, fmt.Errorf("duplicate metadata provided: metadata can only be specified once")
			}
			metadata = v
		default:
			return nil, nil, fmt.Errorf("unsupported argument type: %T", arg)
		}
	}

	return blocks, metadata, nil
}
