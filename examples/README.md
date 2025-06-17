# AI SDK Examples

This directory contains practical examples showing how to use the AI SDK for Go in real applications.

## Prerequisites

Set your API keys as environment variables:

```bash
export OPENAI_API_KEY="your-openai-key-here"
export ANTHROPIC_API_KEY="your-anthropic-key-here"
```

Get your API keys from:
- **OpenAI**: [platform.openai.com/api-keys](https://platform.openai.com/api-keys)
- **Anthropic**: [console.anthropic.com](https://console.anthropic.com/)

## Examples

| Example | Description |
|---------|-------------|
| [**simple-text**](basic/simple-text/) | Generate text from a simple string prompt |
| [**streaming-text**](basic/streaming-text/) | Stream text responses in real-time |

### More Examples Coming Soon

- **Conversation** - Multi-message conversations with context
- **Multi-modal** - Working with images and files
- **Tools** - Function calling and tool usage
- **Advanced** - Production patterns and error handling
- **Real-world** - Complete application examples


## How to Run

Each example is a standalone Go program:

```bash
cd basic/simple-text
go run main.go
```

## Need Help?

- [API Documentation](https://pkg.go.dev/go.jetify.com/ai)
- [Discord Community](https://discord.gg/jetify)
- [GitHub Issues](https://github.com/jetify-com/ai/issues) 