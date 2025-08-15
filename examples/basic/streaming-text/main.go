package main

import (
	"context"
	"fmt"
	"log"

	"go.jetify.com/ai"
	"go.jetify.com/ai/api"
	"go.jetify.com/ai/provider/openai"
)

func example() error {
	// Create a model
	model := openai.NewLanguageModel("gpt-4o-mini")

	// Stream text
	response, err := ai.StreamTextStr(
		context.Background(),
		"Explain what artificial intelligence is in simple terms",
		ai.WithModel(model),
		ai.WithMaxOutputTokens(100),
	)
	if err != nil {
		return err
	}

	// Print the streaming response:
	printStreamResponse(response)

	return nil
}

func printStreamResponse(response *api.StreamResponse) {
	fmt.Print("AI: ")

	for event := range response.Stream {
		switch e := event.(type) {
		case *api.TextDeltaEvent:
			// Print text delta events as they arrive
			fmt.Print(e.TextDelta)
		case *api.FinishEvent:
			// Print final information
			fmt.Printf("\n\nFinish Reason: %s\n", e.FinishReason)
			fmt.Printf("Usage: Input=%d, Output=%d, Total=%d tokens\n",
				e.Usage.InputTokens,
				e.Usage.OutputTokens,
				e.Usage.TotalTokens)
		case *api.ErrorEvent:
			// Handle errors
			fmt.Printf("\nError: %s\n", e.Error())
		}
	}
}

func main() {
	if err := example(); err != nil {
		log.Fatal(err)
	}
}
