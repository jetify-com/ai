# AI SDK for Go

### Build powerful AI applications and agents using a unified API.

[![Version](https://img.shields.io/github/v/release/jetify-com/ai?color=green&label=version&sort=semver)](https://github.com/jetify-com/ai/releases)
[![Go Reference](https://pkg.go.dev/badge/go.jetify.com/ai)](https://pkg.go.dev/go.jetify.com/ai)
[![License](https://img.shields.io/github/license/jetify-com/ai)]()
[![Join Discord](https://img.shields.io/discord/903306922852245526?color=7389D8&label=discord&logo=discord&logoColor=ffffff&cacheSeconds=1800)](https://discord.gg/jetify)

*Primary Author(s)*: [Daniel Loreto](https://github.com/loreto)

## Introduction

Jetify's **AI SDK for Go** is a unified interface for interacting with multiple AI providers including OpenAI, Anthropic, and more.
Inspired by [Vercel's AI SDK](https://github.com/vercel/ai) for TypeScript, we bring a similar developer experience to the Go ecosystem.

It is maintained and developed by [Jetify](https://www.jetify.com). We are in the process of migrating our production code
to use this SDK as the primary way our AI agents integrate with different LLM providers.

### The Problem

Building AI applications go today means dealing with:
- **Fragmented ecosystems** - Each provider has different APIs, authentication, and patterns
- **Vendor lock-in** - Switching providers requires rewriting significant application code
- **Poor Go developer experience** - Official Go SDKs are often auto-generated from OpenAPI specs, resulting in unidiomatic Go code
- **Complex multi-modal handling** - Different providers handle images, files, and tools differently

### Our Solution

The AI SDK provides a **unified interface** across multiple AI providers, with key advantages:

1. **Provider abstraction** - Common interfaces for language models, embeddings, and image generation
2. **Go-first design** - Built specifically for Go developers with idiomatic patterns and strong typing
3. **Production-ready** - Comprehensive error handling, automatic retries, rate limiting, and robust provider failover
4. **Multi-modal by default** - First-class support for text, images, files, and structured outputs across all providers
5. **Extensible architecture** - Clean interfaces make it easy to add new providers while maintaining backward compatibility

## Features

* [x] **Multi-Provider Support** – [OpenAI](#), [Anthropic](#), with more coming
* [x] **Multi-Modal Inputs** – Text, images, and files in conversations
* [x] **Tool Calling** – Function calling with parallel execution
* [x] **Language Models** – Text generation with streaming support
* [ ] **Embedding Models** – Text embeddings for semantic search
* [ ] **Image Models** – Generate images from text prompts
* [ ] **Structured Outputs** – JSON generation with schema validation

### Language Models

* [x] Text generation (streaming & non-streaming)
* [x] Multi-modal conversations (text + images + files)
* [x] System messages and conversation history
* [x] Tool/function calling with structured schemas
* [ ] JSON output with schema validation

### Provider-Specific Features

* [x] **OpenAI** - Web search, computer use, file search tools
* [x] **Anthropic** - Claude's advanced reasoning and tool use

## Status

- [x] Private Alpha: We are testing the SDK with a select group of developers.
- [x] Public Alpha: Open to all developers, but breaking changes still expected.
- [ ] Public Beta: Stable enough for most non-enterprise use cases.
- [ ] General Availability (v1): Ready for production use at scale with guaranteed API stability.

We are currently in **Public Alpha**. The SDK functionality is stable but the API may have breaking changes. While in alpha, minor version bumps indicate breaking changes (`0.1.0` -> `0.2.0` would indicate a breaking change). Watch "releases" of this repo to get notified of major updates.

## Installation

```bash
go get go.jetify.com/ai
```

## Quickstart

Get started with a simple text generation example:

```go
package main

import (
    "context"
    "fmt"
    "log"

    "go.jetify.com/ai"
    "go.jetify.com/ai/provider/openai"
)

func main() {
    // Set up your model
    model := openai.NewLanguageModel("gpt-4o")

    // Generate text
    response, err := ai.GenerateTextStr(
        context.Background(),
        "Explain quantum computing in simple terms",
        ai.WithModel(model),
        ai.WithMaxTokens(200),
    )
    if err != nil {
        log.Fatal(err)
    }

    // Do whatever you want with the response...
}
```

For detailed examples, see our [examples directory](examples/).

## Documentation

Comprehensive documentation is available:

* **[API Reference](https://pkg.go.dev/go.jetify.com/ai)** - Complete Go package documentation
* **[Examples](examples/)** - Real-world usage patterns

## Community & Support

Join our community and get help:

* **Discord** – [https://discord.gg/jetify](https://discord.gg/jetify) (best for quick questions & showcase)
* **GitHub Discussions** – [Discussions](https://github.com/jetify-com/ai/discussions) (best for ideas & design questions)
* **Issues** – [Bug reports & feature requests](https://github.com/jetify-com/ai/issues)

## Contributing

We 💖 contributions! Please read [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

Licensed under the **Apache 2.0 License** – see [LICENSE](LICENSE) for details.
