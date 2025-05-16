We are working on an AI SDK that allows you to use a common interface against different LLM providers.

We've already implemented an openrouter provider, and now we want to implement an Anthropic provider using the anthropic go client  @https://github.com/anthropics/anthropic-sdk-go 

You can use go doc commands to understand the anthropic SDK. The import path is:
```
import (
	"github.com/anthropics/anthropic-sdk-go" // imported as anthropic
)
```

We need to translate our AI SDK prompt types defined in @llm_prompt.go into the corresponding types from the anthropic sdk.

Implement the encoding functions in: @encode_prompt.go 

But first take a look at @encode_prompt.go to see how we did it for the openrouter case (we want to do something analougous)

