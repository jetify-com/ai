package openrouter

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.jetify.com/ai/api"
)

// openRouterErrorData matches the JSON structure of OpenRouter error responses.
type openRouterErrorData struct {
	Error struct {
		Message string  `json:"message"`
		Type    string  `json:"type"`
		Param   any     `json:"param"`
		Code    *string `json:"code"`
	} `json:"error"`
}

// parseOpenRouterErrorJSON attempts to unmarshal the body into openRouterErrorData.
func parseOpenRouterErrorJSON(body []byte) (*openRouterErrorData, error) {
	var parsed openRouterErrorData
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, api.NewJSONParseError(string(body), err)
	}
	return &parsed, nil
}

// OpenRouterFailedResponseHandler constructs an APICallError from a non-2xx OpenRouter response.
func OpenRouterFailedResponseHandler(resp *http.Response, rawBody []byte, requestBody any) error {
	parsed, err := parseOpenRouterErrorJSON(rawBody)
	if err == nil {
		return &api.APICallError{
			AISDKError: api.NewAISDKError("AI_APICallError", parsed.Error.Message, nil),
			URL:        resp.Request.URL,
			Request:    resp.Request,
			StatusCode: resp.StatusCode,
			Response:   resp,
			Data:       parsed,
		}
	}

	// Fallback if we cannot parse the error JSON
	return &api.APICallError{
		AISDKError: api.NewAISDKError("AI_APICallError", fmt.Sprintf("%d %s", resp.StatusCode, resp.Status), err),
		URL:        resp.Request.URL,
		Request:    resp.Request,
		StatusCode: resp.StatusCode,
		Response:   resp,
	}
}
