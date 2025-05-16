package codec

import (
	"go.jetify.com/ai/api"
	"go.jetify.com/ai/provider/openrouter/internal/client"
)

// DecodeLogProbs converts OpenRouter's chat logprobs format to the SDK's LogProb format
func DecodeLogProbs(logprobs *client.LogProbs) []api.LogProb {
	if logprobs == nil || logprobs.Content == nil {
		return []api.LogProb{} // Return empty slice instead of nil
	}

	result := make([]api.LogProb, len(logprobs.Content))
	for i, item := range logprobs.Content {
		result[i] = decodeChatLogProb(item)
	}

	return result
}

// decodeChatLogProb converts a single OpenRouter LogProb to the SDK's LogProb format
func decodeChatLogProb(item client.LogProb) api.LogProb {
	return api.LogProb{
		Token:       item.Token,
		LogProb:     item.LogProb,
		TopLogProbs: decodeTokenLogProbs(item.TopLogProbs),
	}
}

// DecodeCompletionLogProbs converts OpenRouter's completion logprobs format to the SDK's LogProb format.
// It handles nil input by returning an empty slice.
func DecodeCompletionLogProbs(logprobs *client.CompletionLogProbs) []api.LogProb {
	if logprobs == nil {
		return []api.LogProb{}
	}

	result := make([]api.LogProb, len(logprobs.Tokens))
	for i, token := range logprobs.Tokens {
		result[i] = decodeCompletionLogProbForToken(logprobs, i, token)
	}

	return result
}

// decodeCompletionLogProbForToken converts a single token's logprob data into the SDK's LogProb format
func decodeCompletionLogProbForToken(logprobs *client.CompletionLogProbs, index int, token string) api.LogProb {
	return api.LogProb{
		Token:       token,
		LogProb:     getCompletionLogProbValue(logprobs.TokenLogProbs, index),
		TopLogProbs: decodeCompletionTopLogProbs(logprobs.TopLogProbs, index),
	}
}

// getCompletionLogProbValue safely gets the logprob value for a token, defaulting to 0 if out of bounds
func getCompletionLogProbValue(logprobs []float64, index int) float64 {
	if index < len(logprobs) {
		return logprobs[index]
	}
	return 0
}

// decodeCompletionTopLogProbs converts the top logprobs map for a token into the SDK's TokenLogProb format
func decodeCompletionTopLogProbs(topLogProbs []map[string]float64, index int) []api.TokenLogProb {
	if topLogProbs == nil || index >= len(topLogProbs) || topLogProbs[index] == nil {
		return []api.TokenLogProb{} // Always return empty slice instead of nil
	}

	topMap := topLogProbs[index]
	pairs := make([]client.TopLogProb, 0, len(topMap))

	for token, logprob := range topMap {
		pairs = append(pairs, client.TopLogProb{
			Token:   token,
			LogProb: logprob,
		})
	}

	return decodeTokenLogProbs(pairs)
}

// decodeTokenLogProbs converts a slice of token/logprob pairs to the SDK's TokenLogProb format
func decodeTokenLogProbs(pairs []client.TopLogProb) []api.TokenLogProb {
	result := make([]api.TokenLogProb, len(pairs))
	for i, pair := range pairs {
		result[i] = api.TokenLogProb{
			Token:   pair.Token,
			LogProb: pair.LogProb,
		}
	}
	return result
}
