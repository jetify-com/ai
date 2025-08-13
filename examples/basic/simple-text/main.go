package main

import (
	"context"
	"log"

	"github.com/k0kubun/pp/v3"
	"go.jetify.com/ai"
	"go.jetify.com/ai/api"
	"go.jetify.com/ai/provider/openai"
)

func example() error {
	// Create a model
	model := openai.NewLanguageModel("gpt-4o-mini")

	// Generate text
	response, err := ai.GenerateTextStr(
		context.Background(),
		"Explain what artificial intelligence is in simple terms",
		ai.WithModel(model),
		ai.WithMaxOutputTokens(100),
	)
	if err != nil {
		return err
	}

	// Print the response:
	printResponse(response)

	return nil
}

func printResponse(response api.Response) {
	response.ProviderMetadata = nil
	response.Warnings = nil
	printer := pp.New()
	printer.SetOmitEmpty(true)
	_, _ = printer.Print(response)
}

func main() {
	if err := example(); err != nil {
		log.Fatal(err)
	}
}
