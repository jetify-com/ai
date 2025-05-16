package model

// Last updated: 2025-02-13

// Dolphin30R1Mistral24bFree is the ID for model Dolphin3.0 R1 Mistral 24b (free)
//
// Dolphin 3.0 R1 is the next generation of the Dolphin series of instruct-tuned models.  Designed to be the ultimate general purpose local model, enabling coding, math, agentic, function calling, and general use cases.
//
// The R1 version has been trained for 3 epochs to reason using 800k reasoning traces from the Dolphin-R1 dataset.
//
// Dolphin aims to be a general purpose reasoning instruct model, similar to the models behind ChatGPT, Claude, Gemini.
//
// Part of the [Dolphin 3.0 Collection](https://huggingface.co/collections/cognitivecomputations/dolphin-30-677ab47f73d7ff66743979a3) Curated and trained by [Eric Hartford](https://huggingface.co/ehartford), [Ben Gitter](https://huggingface.co/bigstorm), [BlouseJury](https://huggingface.co/BlouseJury) and [Cognitive Computations](https://huggingface.co/cognitivecomputations)
const Dolphin30R1Mistral24bFree = "cognitivecomputations/dolphin3.0-r1-mistral-24b:free"

// Dolphin30Mistral24bFree is the ID for model Dolphin3.0 Mistral 24b (free)
//
// Dolphin 3.0 is the next generation of the Dolphin series of instruct-tuned models.  Designed to be the ultimate general purpose local model, enabling coding, math, agentic, function calling, and general use cases.
//
// Dolphin aims to be a general purpose instruct model, similar to the models behind ChatGPT, Claude, Gemini.
//
// Part of the [Dolphin 3.0 Collection](https://huggingface.co/collections/cognitivecomputations/dolphin-30-677ab47f73d7ff66743979a3) Curated and trained by [Eric Hartford](https://huggingface.co/ehartford), [Ben Gitter](https://huggingface.co/bigstorm), [BlouseJury](https://huggingface.co/BlouseJury) and [Cognitive Computations](https://huggingface.co/cognitivecomputations)
const Dolphin30Mistral24bFree = "cognitivecomputations/dolphin3.0-mistral-24b:free"

// LlamaGuard38b is the ID for model Llama Guard 3 8b
//
// Llama Guard 3 is a Llama-3.1-8B pretrained model, fine-tuned for content safety classification. Similar to previous versions, it can be used to classify content in both LLM inputs (prompt classification) and in LLM responses (response classification). It acts as an LLM – it generates text in its output that indicates whether a given prompt or response is safe or unsafe, and if unsafe, it also lists the content categories violated.
//
// Llama Guard 3 was aligned to safeguard against the MLCommons standardized hazards taxonomy and designed to support Llama 3.1 capabilities. Specifically, it provides content moderation in 8 languages, and was optimized to support safety and security for search and code interpreter tool calls.
const LlamaGuard38b = "meta-llama/llama-guard-3-8b"

// OpenAIO3MiniHigh is the ID for model OpenAI: o3 Mini High
//
// OpenAI o3-mini-high is the same model as [o3-mini](/openai/o3-mini) with reasoning_effort set to high.
//
// o3-mini is a cost-efficient language model optimized for STEM reasoning tasks, particularly excelling in science, mathematics, and coding. The model features three adjustable reasoning effort levels and supports key developer capabilities including function calling, structured outputs, and streaming, though it does not include vision processing capabilities.
//
// The model demonstrates significant improvements over its predecessor, with expert testers preferring its responses 56% of the time and noting a 39% reduction in major errors on complex questions. With medium reasoning effort settings, o3-mini matches the performance of the larger o1 model on challenging reasoning evaluations like AIME and GPQA, while maintaining lower latency and cost.
const OpenAIO3MiniHigh = "openai/o3-mini-high"

// Llama31Tulu3405b is the ID for model Llama 3.1 Tulu 3 405b
//
// Tülu 3 405B is the largest model in the Tülu 3 family, applying fully open post-training recipes at a 405B parameter scale. Built on the Llama 3.1 405B base, it leverages Reinforcement Learning with Verifiable Rewards (RLVR) to enhance instruction following, MATH, GSM8K, and IFEval performance. As part of Tülu 3’s fully open-source approach, it offers state-of-the-art capabilities while surpassing prior open-weight models like Llama 3.1 405B Instruct and Nous Hermes 3 405B on multiple benchmarks. To read more, [click here.](https://allenai.org/blog/tulu-3-405B)
const Llama31Tulu3405b = "allenai/llama-3.1-tulu-3-405b"

// DeepSeekR1DistillLlama8B is the ID for model DeepSeek: R1 Distill Llama 8B
//
// DeepSeek R1 Distill Llama 8B is a distilled large language model based on [Llama-3.1-8B-Instruct](/meta-llama/llama-3.1-8b-instruct), using outputs from [DeepSeek R1](/deepseek/deepseek-r1). The model combines advanced distillation techniques to achieve high performance across multiple benchmarks, including:
//
// - AIME 2024 pass@1: 50.4
// - MATH-500 pass@1: 89.1
// - CodeForces Rating: 1205
//
// The model leverages fine-tuning from DeepSeek R1's outputs, enabling competitive performance comparable to larger frontier models.
//
// Hugging Face:
// - [Llama-3.1-8B](https://huggingface.co/meta-llama/Llama-3.1-8B)
// - [DeepSeek-R1-Distill-Llama-8B](https://huggingface.co/deepseek-ai/DeepSeek-R1-Distill-Llama-8B)   |
const DeepSeekR1DistillLlama8B = "deepseek/deepseek-r1-distill-llama-8b"

// GoogleGeminiFlash20 is the ID for model Google: Gemini Flash 2.0
//
// Gemini Flash 2.0 offers a significantly faster time to first token (TTFT) compared to [Gemini Flash 1.5](/google/gemini-flash-1.5), while maintaining quality on par with larger models like [Gemini Pro 1.5](/google/gemini-pro-1.5). It introduces notable enhancements in multimodal understanding, coding capabilities, complex instruction following, and function calling. These advancements come together to deliver more seamless and robust agentic experiences.
const GoogleGeminiFlash20 = "google/gemini-2.0-flash-001"

// GoogleGeminiFlashLite20PreviewFree is the ID for model Google: Gemini Flash Lite 2.0 Preview (free)
//
// Gemini Flash Lite 2.0 offers a significantly faster time to first token (TTFT) compared to [Gemini Flash 1.5](google/gemini-flash-1.5), while maintaining quality on par with larger models like [Gemini Pro 1.5](google/gemini-pro-1.5). Because it's currently in preview, it will be **heavily rate-limited** by Google. This model will move from free to paid pending a general rollout on February 24th, at $0.075 / $0.30 per million input / output tokens respectively.
const GoogleGeminiFlashLite20PreviewFree = "google/gemini-2.0-flash-lite-preview-02-05:free"

// GoogleGeminiPro20ExperimentalFree is the ID for model Google: Gemini Pro 2.0 Experimental (free)
//
// Gemini 2.0 Pro Experimental is a bleeding-edge version of the Gemini 2.0 Pro model. Because it's currently experimental, it will be **heavily rate-limited** by Google.
//
// Usage of Gemini is subject to Google's [Gemini Terms of Use](https://ai.google.dev/terms).
//
// #multimodal
const GoogleGeminiPro20ExperimentalFree = "google/gemini-2.0-pro-exp-02-05:free"

// QwenQwenVLPlusFree is the ID for model Qwen: Qwen VL Plus (free)
//
// Qwen's Enhanced Large Visual Language Model. Significantly upgraded for detailed recognition capabilities and text recognition abilities, supporting ultra-high pixel resolutions up to millions of pixels and extreme aspect ratios for image input. It delivers significant performance across a broad range of visual tasks.
const QwenQwenVLPlusFree = "qwen/qwen-vl-plus:free"

// AionLabsAion10 is the ID for model AionLabs: Aion-1.0
//
// Aion-1.0 is a multi-model system designed for high performance across various tasks, including reasoning and coding. It is built on DeepSeek-R1, augmented with additional models and techniques such as Tree of Thoughts (ToT) and Mixture of Experts (MoE). It is Aion Lab's most powerful reasoning model.
const AionLabsAion10 = "aion-labs/aion-1.0"

// AionLabsAion10Mini is the ID for model AionLabs: Aion-1.0-Mini
//
// Aion-1.0-Mini 32B parameter model is a distilled version of the DeepSeek-R1 model, designed for strong performance in reasoning domains such as mathematics, coding, and logic. It is a modified variant of a FuseAI model that outperforms R1-Distill-Qwen-32B and R1-Distill-Llama-70B, with benchmark results available on its [Hugging Face page](https://huggingface.co/FuseAI/FuseO1-DeepSeekR1-QwQ-SkyT1-32B-Preview), independently replicated for verification.
const AionLabsAion10Mini = "aion-labs/aion-1.0-mini"

// AionLabsAionRP108B is the ID for model AionLabs: Aion-RP 1.0 (8B)
//
// Aion-RP-Llama-3.1-8B ranks the highest in the character evaluation portion of the RPBench-Auto benchmark, a roleplaying-specific variant of Arena-Hard-Auto, where LLMs evaluate each other’s responses. It is a fine-tuned base model rather than an instruct model, designed to produce more natural and varied writing.
const AionLabsAionRP108B = "aion-labs/aion-rp-llama-3.1-8b"

// QwenQwenTurbo is the ID for model Qwen: Qwen-Turbo
//
// Qwen-Turbo, based on Qwen2.5, is a 1M context model that provides fast speed and low cost, suitable for simple tasks.
const QwenQwenTurbo = "qwen/qwen-turbo"

// QwenQwen25VL72BInstructFree is the ID for model Qwen: Qwen2.5 VL 72B Instruct (free)
//
// Qwen2.5-VL is proficient in recognizing common objects such as flowers, birds, fish, and insects. It is also highly capable of analyzing texts, charts, icons, graphics, and layouts within images.
const QwenQwen25VL72BInstructFree = "qwen/qwen2.5-vl-72b-instruct:free"

// QwenQwenPlus is the ID for model Qwen: Qwen-Plus
//
// Qwen-Plus, based on the Qwen2.5 foundation model, is a 131K context model with a balanced performance, speed, and cost combination.
const QwenQwenPlus = "qwen/qwen-plus"

// QwenQwenMax is the ID for model Qwen: Qwen-Max
//
// Qwen-Max, based on Qwen2.5, provides the best inference performance among [Qwen models](/qwen), especially for complex multi-step tasks. It's a large-scale MoE model that has been pretrained on over 20 trillion tokens and further post-trained with curated Supervised Fine-Tuning (SFT) and Reinforcement Learning from Human Feedback (RLHF) methodologies. The parameter count is unknown.
const QwenQwenMax = "qwen/qwen-max"

// OpenAIO3Mini is the ID for model OpenAI: o3 Mini
//
// OpenAI o3-mini is a cost-efficient language model optimized for STEM reasoning tasks, particularly excelling in science, mathematics, and coding.
//
// This model supports the `reasoning_effort` parameter, which can be set to "high", "medium", or "low" to control the thinking time of the model. The default is "medium". OpenRouter also offers the model slug `openai/o3-mini-high` to default the parameter to "high".
//
// The model features three adjustable reasoning effort levels and supports key developer capabilities including function calling, structured outputs, and streaming, though it does not include vision processing capabilities.
//
// The model demonstrates significant improvements over its predecessor, with expert testers preferring its responses 56% of the time and noting a 39% reduction in major errors on complex questions. With medium reasoning effort settings, o3-mini matches the performance of the larger o1 model on challenging reasoning evaluations like AIME and GPQA, while maintaining lower latency and cost.
const OpenAIO3Mini = "openai/o3-mini"

// DeepSeekR1DistillQwen15B is the ID for model DeepSeek: R1 Distill Qwen 1.5B
//
// DeepSeek R1 Distill Qwen 1.5B is a distilled large language model based on  [Qwen 2.5 Math 1.5B](https://huggingface.co/Qwen/Qwen2.5-Math-1.5B), using outputs from [DeepSeek R1](/deepseek/deepseek-r1). It's a very small and efficient model which outperforms [GPT 4o 0513](/openai/gpt-4o-2024-05-13) on Math Benchmarks.
//
// Other benchmark results include:
//
// - AIME 2024 pass@1: 28.9
// - AIME 2024 cons@64: 52.7
// - MATH-500 pass@1: 83.9
//
// The model leverages fine-tuning from DeepSeek R1's outputs, enabling competitive performance comparable to larger frontier models.
const DeepSeekR1DistillQwen15B = "deepseek/deepseek-r1-distill-qwen-1.5b"

// MistralMistralSmall3Free is the ID for model Mistral: Mistral Small 3 (free)
//
// Mistral Small 3 is a 24B-parameter language model optimized for low-latency performance across common AI tasks. Released under the Apache 2.0 license, it features both pre-trained and instruction-tuned versions designed for efficient local deployment.
//
// The model achieves 81% accuracy on the MMLU benchmark and performs competitively with larger models like Llama 3.3 70B and Qwen 32B, while operating at three times the speed on equivalent hardware. [Read the blog post about the model here.](https://mistral.ai/news/mistral-small-3/)
const MistralMistralSmall3Free = "mistralai/mistral-small-24b-instruct-2501:free"

// MistralMistralSmall3 is the ID for model Mistral: Mistral Small 3
//
// Mistral Small 3 is a 24B-parameter language model optimized for low-latency performance across common AI tasks. Released under the Apache 2.0 license, it features both pre-trained and instruction-tuned versions designed for efficient local deployment.
//
// The model achieves 81% accuracy on the MMLU benchmark and performs competitively with larger models like Llama 3.3 70B and Qwen 32B, while operating at three times the speed on equivalent hardware. [Read the blog post about the model here.](https://mistral.ai/news/mistral-small-3/)
const MistralMistralSmall3 = "mistralai/mistral-small-24b-instruct-2501"

// DeepSeekR1DistillQwen32B is the ID for model DeepSeek: R1 Distill Qwen 32B
//
// DeepSeek R1 Distill Qwen 32B is a distilled large language model based on [Qwen 2.5 32B](https://huggingface.co/Qwen/Qwen2.5-32B), using outputs from [DeepSeek R1](/deepseek/deepseek-r1). It outperforms OpenAI's o1-mini across various benchmarks, achieving new state-of-the-art results for dense models.
//
// Other benchmark results include:
//
// - AIME 2024 pass@1: 72.6
// - MATH-500 pass@1: 94.3
// - CodeForces Rating: 1691
//
// The model leverages fine-tuning from DeepSeek R1's outputs, enabling competitive performance comparable to larger frontier models.
const DeepSeekR1DistillQwen32B = "deepseek/deepseek-r1-distill-qwen-32b"

// DeepSeekR1DistillQwen14B is the ID for model DeepSeek: R1 Distill Qwen 14B
//
// DeepSeek R1 Distill Qwen 14B is a distilled large language model based on [Qwen 2.5 14B](https://huggingface.co/deepseek-ai/DeepSeek-R1-Distill-Qwen-14B), using outputs from [DeepSeek R1](/deepseek/deepseek-r1). It outperforms OpenAI's o1-mini across various benchmarks, achieving new state-of-the-art results for dense models.
//
// Other benchmark results include:
//
// - AIME 2024 pass@1: 69.7
// - MATH-500 pass@1: 93.9
// - CodeForces Rating: 1481
//
// The model leverages fine-tuning from DeepSeek R1's outputs, enabling competitive performance comparable to larger frontier models.
const DeepSeekR1DistillQwen14B = "deepseek/deepseek-r1-distill-qwen-14b"

// PerplexitySonarReasoning is the ID for model Perplexity: Sonar Reasoning
//
// Sonar Reasoning is a reasoning model provided by Perplexity based on [DeepSeek R1](/deepseek/deepseek-r1).
//
// It allows developers to utilize long chain of thought with built-in web search. Sonar Reasoning is uncensored and hosted in US datacenters.
const PerplexitySonarReasoning = "perplexity/sonar-reasoning"

// PerplexitySonar is the ID for model Perplexity: Sonar
//
// Sonar is lightweight, affordable, fast, and simple to use — now featuring citations and the ability to customize sources. It is designed for companies seeking to integrate lightweight question-and-answer features optimized for speed.
const PerplexitySonar = "perplexity/sonar"

// LiquidLFM7B is the ID for model Liquid: LFM 7B
//
// LFM-7B, a new best-in-class language model. LFM-7B is designed for exceptional chat capabilities, including languages like Arabic and Japanese. Powered by the Liquid Foundation Model (LFM) architecture, it exhibits unique features like low memory footprint and fast inference speed.
//
// LFM-7B is the world’s best-in-class multilingual language model in English, Arabic, and Japanese.
//
// See the [launch announcement](https://www.liquid.ai/lfm-7b) for benchmarks and more info.
const LiquidLFM7B = "liquid/lfm-7b"

// LiquidLFM3B is the ID for model Liquid: LFM 3B
//
// Liquid's LFM 3B delivers incredible performance for its size. It positions itself as first place among 3B parameter transformers, hybrids, and RNN models It is also on par with Phi-3.5-mini on multiple benchmarks, while being 18.4% smaller.
//
// LFM-3B is the ideal choice for mobile and other edge text-based applications.
//
// See the [launch announcement](https://www.liquid.ai/liquid-foundation-models) for benchmarks and more info.
const LiquidLFM3B = "liquid/lfm-3b"

// DeepSeekR1DistillLlama70BFree is the ID for model DeepSeek: R1 Distill Llama 70B (free)
//
// DeepSeek R1 Distill Llama 70B is a distilled large language model based on [Llama-3.3-70B-Instruct](/meta-llama/llama-3.3-70b-instruct), using outputs from [DeepSeek R1](/deepseek/deepseek-r1). The model combines advanced distillation techniques to achieve high performance across multiple benchmarks, including:
//
// - AIME 2024 pass@1: 70.0
// - MATH-500 pass@1: 94.5
// - CodeForces Rating: 1633
//
// The model leverages fine-tuning from DeepSeek R1's outputs, enabling competitive performance comparable to larger frontier models.
const DeepSeekR1DistillLlama70BFree = "deepseek/deepseek-r1-distill-llama-70b:free"

// DeepSeekR1DistillLlama70B is the ID for model DeepSeek: R1 Distill Llama 70B
//
// DeepSeek R1 Distill Llama 70B is a distilled large language model based on [Llama-3.3-70B-Instruct](/meta-llama/llama-3.3-70b-instruct), using outputs from [DeepSeek R1](/deepseek/deepseek-r1). The model combines advanced distillation techniques to achieve high performance across multiple benchmarks, including:
//
// - AIME 2024 pass@1: 70.0
// - MATH-500 pass@1: 94.5
// - CodeForces Rating: 1633
//
// The model leverages fine-tuning from DeepSeek R1's outputs, enabling competitive performance comparable to larger frontier models.
const DeepSeekR1DistillLlama70B = "deepseek/deepseek-r1-distill-llama-70b"

// GoogleGemini20FlashThinkingExperimental0121Free is the ID for model Google: Gemini 2.0 Flash Thinking Experimental 01-21 (free)
//
// Gemini 2.0 Flash Thinking Experimental (01-21) is a snapshot of Gemini 2.0 Flash Thinking Experimental.
//
// Gemini 2.0 Flash Thinking Mode is an experimental model that's trained to generate the "thinking process" the model goes through as part of its response. As a result, Thinking Mode is capable of stronger reasoning capabilities in its responses than the [base Gemini 2.0 Flash model](/google/gemini-2.0-flash-exp).
const GoogleGemini20FlashThinkingExperimental0121Free = "google/gemini-2.0-flash-thinking-exp:free"

// DeepSeekR1Free is the ID for model DeepSeek: R1 (free)
//
// DeepSeek R1 is here: Performance on par with [OpenAI o1](/openai/o1), but open-sourced and with fully open reasoning tokens. It's 671B parameters in size, with 37B active in an inference pass.
//
// Fully open-source model & [technical report](https://api-docs.deepseek.com/news/news250120).
//
// MIT licensed: Distill & commercialize freely!
const DeepSeekR1Free = "deepseek/deepseek-r1:free"

// DeepSeekR1 is the ID for model DeepSeek: R1
//
// DeepSeek R1 is here: Performance on par with [OpenAI o1](/openai/o1), but open-sourced and with fully open reasoning tokens. It's 671B parameters in size, with 37B active in an inference pass.
//
// Fully open-source model & [technical report](https://api-docs.deepseek.com/news/news250120).
//
// MIT licensed: Distill & commercialize freely!
const DeepSeekR1 = "deepseek/deepseek-r1"

// RogueRose103BV02Free is the ID for model Rogue Rose 103B v0.2 (free)
//
// Rogue Rose demonstrates strong capabilities in roleplaying and storytelling applications, potentially surpassing other models in the 103-120B parameter range. While it occasionally exhibits inconsistencies with scene logic, the overall interaction quality represents an advancement in natural language processing for creative applications.
//
// It is a 120-layer frankenmerge model combining two custom 70B architectures from November 2023, derived from the [xwin-stellarbright-erp-70b-v2](https://huggingface.co/sophosympatheia/xwin-stellarbright-erp-70b-v2) base.
const RogueRose103BV02Free = "sophosympatheia/rogue-rose-103b-v0.2:free"

// MiniMaxMiniMax01 is the ID for model MiniMax: MiniMax-01
//
// MiniMax-01 is a combines MiniMax-Text-01 for text generation and MiniMax-VL-01 for image understanding. It has 456 billion parameters, with 45.9 billion parameters activated per inference, and can handle a context of up to 4 million tokens.
//
// The text model adopts a hybrid architecture that combines Lightning Attention, Softmax Attention, and Mixture-of-Experts (MoE). The image model adopts the “ViT-MLP-LLM” framework and is trained on top of the text model.
//
// To read more about the release, see: https://www.minimaxi.com/en/news/minimax-01-series-2
const MiniMaxMiniMax01 = "minimax/minimax-01"

// MistralCodestral2501 is the ID for model Mistral: Codestral 2501
//
// [Mistral](/mistralai)'s cutting-edge language model for coding. Codestral specializes in low-latency, high-frequency tasks such as fill-in-the-middle (FIM), code correction and test generation.
//
// Learn more on their blog post: https://mistral.ai/news/codestral-2501/
const MistralCodestral2501 = "mistralai/codestral-2501"

// MicrosoftPhi4 is the ID for model Microsoft: Phi 4
//
// [Microsoft Research](/microsoft) Phi-4 is designed to perform well in complex reasoning tasks and can operate efficiently in situations with limited memory or where quick responses are needed.
//
// At 14 billion parameters, it was trained on a mix of high-quality synthetic datasets, data from curated websites, and academic materials. It has undergone careful improvement to follow instructions accurately and maintain strong safety standards. It works best with English language inputs.
//
// For more information, please see [Phi-4 Technical Report](https://arxiv.org/pdf/2412.08905)
const MicrosoftPhi4 = "microsoft/phi-4"

// Sao10KLlama3170BHanamiX1 is the ID for model Sao10K: Llama 3.1 70B Hanami x1
//
// This is [Sao10K](/sao10k)'s experiment over [Euryale v2.2](/sao10k/l3.1-euryale-70b).
const Sao10KLlama3170BHanamiX1 = "sao10k/l3.1-70b-hanami-x1"

// DeepSeekDeepSeekV3Free is the ID for model DeepSeek: DeepSeek V3 (free)
//
// DeepSeek-V3 is the latest model from the DeepSeek team, building upon the instruction following and coding abilities of the previous versions. Pre-trained on nearly 15 trillion tokens, the reported evaluations reveal that the model outperforms other open-source models and rivals leading closed-source models.
//
// For model details, please visit [the DeepSeek-V3 repo](https://github.com/deepseek-ai/DeepSeek-V3) for more information, or see the [launch announcement](https://api-docs.deepseek.com/news/news1226).
const DeepSeekDeepSeekV3Free = "deepseek/deepseek-chat:free"

// DeepSeekDeepSeekV3 is the ID for model DeepSeek: DeepSeek V3
//
// DeepSeek-V3 is the latest model from the DeepSeek team, building upon the instruction following and coding abilities of the previous versions. Pre-trained on nearly 15 trillion tokens, the reported evaluations reveal that the model outperforms other open-source models and rivals leading closed-source models.
//
// For model details, please visit [the DeepSeek-V3 repo](https://github.com/deepseek-ai/DeepSeek-V3) for more information, or see the [launch announcement](https://api-docs.deepseek.com/news/news1226).
const DeepSeekDeepSeekV3 = "deepseek/deepseek-chat"

// QwenQvQ72BPreview is the ID for model Qwen: QvQ 72B Preview
//
// QVQ-72B-Preview is an experimental research model developed by the [Qwen](/qwen) team, focusing on enhancing visual reasoning capabilities.
//
// ## Performance
//
// |                | **QVQ-72B-Preview** | o1-2024-12-17 | gpt-4o-2024-05-13 | Claude3.5 Sonnet-20241022 | Qwen2VL-72B |
// |----------------|-----------------|---------------|-------------------|----------------------------|-------------|
// | MMMU(val)      | 70.3            | 77.3          | 69.1              | 70.4                       | 64.5        |
// | MathVista(mini) | 71.4            | 71.0          | 63.8              | 65.3                       | 70.5        |
// | MathVision(full)   | 35.9            | –             | 30.4              | 35.6                       | 25.9        |
// | OlympiadBench  | 20.4            | –             | 25.9              | –                          | 11.2        |
//
// ## Limitations
//
// 1. **Language Mixing and Code-Switching:** The model might occasionally mix different languages or unexpectedly switch between them, potentially affecting the clarity of its responses.
// 2. **Recursive Reasoning Loops:**  There's a risk of the model getting caught in recursive reasoning loops, leading to lengthy responses that may not even arrive at a final answer.
// 3. **Safety and Ethical Considerations:** Robust safety measures are needed to ensure reliable and safe performance. Users should exercise caution when deploying this model.
// 4. **Performance and Benchmark Limitations:** Despite the improvements in visual reasoning, QVQ doesn’t entirely replace the capabilities of [Qwen2-VL-72B](/qwen/qwen-2-vl-72b-instruct). During multi-step visual reasoning, the model might gradually lose focus on the image content, leading to hallucinations. Moreover, QVQ doesn’t show significant improvement over [Qwen2-VL-72B](/qwen/qwen-2-vl-72b-instruct) in basic recognition tasks like identifying people, animals, or plants.
//
// Note: Currently, the model only supports single-round dialogues and image outputs. It does not support video inputs.
const QwenQvQ72BPreview = "qwen/qvq-72b-preview"

// GoogleGemini20FlashThinkingExperimentalFree is the ID for model Google: Gemini 2.0 Flash Thinking Experimental (free)
//
// Gemini 2.0 Flash Thinking Mode is an experimental model that's trained to generate the "thinking process" the model goes through as part of its response. As a result, Thinking Mode is capable of stronger reasoning capabilities in its responses than the [base Gemini 2.0 Flash model](/google/gemini-2.0-flash-exp).
const GoogleGemini20FlashThinkingExperimentalFree = "google/gemini-2.0-flash-thinking-exp-1219:free"

// Sao10KLlama33Euryale70B is the ID for model Sao10K: Llama 3.3 Euryale 70B
//
// Euryale L3.3 70B is a model focused on creative roleplay from [Sao10k](https://ko-fi.com/sao10k). It is the successor of [Euryale L3 70B v2.2](/models/sao10k/l3-euryale-70b).
const Sao10KLlama33Euryale70B = "sao10k/l3.3-euryale-70b"

// OpenAIO1 is the ID for model OpenAI: o1
//
// The latest and strongest model family from OpenAI, o1 is designed to spend more time thinking before responding. The o1 model series is trained with large-scale reinforcement learning to reason using chain of thought.
//
// The o1 models are optimized for math, science, programming, and other STEM-related tasks. They consistently exhibit PhD-level accuracy on benchmarks in physics, chemistry, and biology. Learn more in the [launch announcement](https://openai.com/o1).
const OpenAIO1 = "openai/o1"

// EVALlama33370b is the ID for model EVA Llama 3.33 70b
//
// EVA Llama 3.33 70b is a roleplay and storywriting specialist model. It is a full-parameter finetune of [Llama-3.3-70B-Instruct](https://openrouter.ai/meta-llama/llama-3.3-70b-instruct) on mixture of synthetic and natural data.
//
// It uses Celeste 70B 0.1 data mixture, greatly expanding it to improve versatility, creativity and "flavor" of the resulting model
//
// This model was built with Llama by Meta.
const EVALlama33370b = "eva-unit-01/eva-llama-3.33-70b"

// XAIGrok2Vision1212 is the ID for model xAI: Grok 2 Vision 1212
//
// Grok 2 Vision 1212 advances image-based AI with stronger visual comprehension, refined instruction-following, and multilingual support. From object recognition to style analysis, it empowers developers to build more intuitive, visually aware applications. Its enhanced steerability and reasoning establish a robust foundation for next-generation image solutions.
//
// To read more about this model, check out [xAI's announcement](https://x.ai/blog/grok-1212).
const XAIGrok2Vision1212 = "x-ai/grok-2-vision-1212"

// XAIGrok21212 is the ID for model xAI: Grok 2 1212
//
// Grok 2 1212 introduces significant enhancements to accuracy, instruction adherence, and multilingual support, making it a powerful and flexible choice for developers seeking a highly steerable, intelligent model.
const XAIGrok21212 = "x-ai/grok-2-1212"

// CohereCommandR7B122024 is the ID for model Cohere: Command R7B (12-2024)
//
// Command R7B (12-2024) is a small, fast update of the Command R+ model, delivered in December 2024. It excels at RAG, tool use, agents, and similar tasks requiring complex reasoning and multiple steps.
//
// Use of this model is subject to Cohere's [Usage Policy](https://docs.cohere.com/docs/usage-policy) and [SaaS Agreement](https://cohere.com/saas-agreement).
const CohereCommandR7B122024 = "cohere/command-r7b-12-2024"

// GoogleGeminiFlash20ExperimentalFree is the ID for model Google: Gemini Flash 2.0 Experimental (free)
//
// Gemini Flash 2.0 offers a significantly faster time to first token (TTFT) compared to [Gemini Flash 1.5](/google/gemini-flash-1.5), while maintaining quality on par with larger models like [Gemini Pro 1.5](/google/gemini-pro-1.5). It introduces notable enhancements in multimodal understanding, coding capabilities, complex instruction following, and function calling. These advancements come together to deliver more seamless and robust agentic experiences.
const GoogleGeminiFlash20ExperimentalFree = "google/gemini-2.0-flash-exp:free"

// GoogleGeminiExperimental1206Free is the ID for model Google: Gemini Experimental 1206 (free)
//
// Experimental release (December 6, 2024) of Gemini.
const GoogleGeminiExperimental1206Free = "google/gemini-exp-1206:free"

// MetaLlama3370BInstructFree is the ID for model Meta: Llama 3.3 70B Instruct (free)
//
// The Meta Llama 3.3 multilingual large language model (LLM) is a pretrained and instruction tuned generative model in 70B (text in/text out). The Llama 3.3 instruction tuned text only model is optimized for multilingual dialogue use cases and outperforms many of the available open source and closed chat models on common industry benchmarks.
//
// Supported languages: English, German, French, Italian, Portuguese, Hindi, Spanish, and Thai.
//
// [Model Card](https://github.com/meta-llama/llama-models/blob/main/models/llama3_3/CARD.md)
const MetaLlama3370BInstructFree = "meta-llama/llama-3.3-70b-instruct:free"

// MetaLlama3370BInstruct is the ID for model Meta: Llama 3.3 70B Instruct
//
// The Meta Llama 3.3 multilingual large language model (LLM) is a pretrained and instruction tuned generative model in 70B (text in/text out). The Llama 3.3 instruction tuned text only model is optimized for multilingual dialogue use cases and outperforms many of the available open source and closed chat models on common industry benchmarks.
//
// Supported languages: English, German, French, Italian, Portuguese, Hindi, Spanish, and Thai.
//
// [Model Card](https://github.com/meta-llama/llama-models/blob/main/models/llama3_3/CARD.md)
const MetaLlama3370BInstruct = "meta-llama/llama-3.3-70b-instruct"

// AmazonNovaLite10 is the ID for model Amazon: Nova Lite 1.0
//
// Amazon Nova Lite 1.0 is a very low-cost multimodal model from Amazon that focused on fast processing of image, video, and text inputs to generate text output. Amazon Nova Lite can handle real-time customer interactions, document analysis, and visual question-answering tasks with high accuracy.
//
// With an input context of 300K tokens, it can analyze multiple images or up to 30 minutes of video in a single input.
const AmazonNovaLite10 = "amazon/nova-lite-v1"

// AmazonNovaMicro10 is the ID for model Amazon: Nova Micro 1.0
//
// Amazon Nova Micro 1.0 is a text-only model that delivers the lowest latency responses in the Amazon Nova family of models at a very low cost. With a context length of 128K tokens and optimized for speed and cost, Amazon Nova Micro excels at tasks such as text summarization, translation, content classification, interactive chat, and brainstorming. It has  simple mathematical reasoning and coding abilities.
const AmazonNovaMicro10 = "amazon/nova-micro-v1"

// AmazonNovaPro10 is the ID for model Amazon: Nova Pro 1.0
//
// Amazon Nova Pro 1.0 is a capable multimodal model from Amazon focused on providing a combination of accuracy, speed, and cost for a wide range of tasks. As of December 2024, it achieves state-of-the-art performance on key benchmarks including visual question answering (TextVQA) and video understanding (VATEX).
//
// Amazon Nova Pro demonstrates strong capabilities in processing both visual and textual information and at analyzing financial documents.
//
// **NOTE**: Video input is not supported at this time.
const AmazonNovaPro10 = "amazon/nova-pro-v1"

// QwenQwQ32BPreview is the ID for model Qwen: QwQ 32B Preview
//
// QwQ-32B-Preview is an experimental research model focused on AI reasoning capabilities developed by the Qwen Team. As a preview release, it demonstrates promising analytical abilities while having several important limitations:
//
// 1. **Language Mixing and Code-Switching**: The model may mix languages or switch between them unexpectedly, affecting response clarity.
// 2. **Recursive Reasoning Loops**: The model may enter circular reasoning patterns, leading to lengthy responses without a conclusive answer.
// 3. **Safety and Ethical Considerations**: The model requires enhanced safety measures to ensure reliable and secure performance, and users should exercise caution when deploying it.
// 4. **Performance and Benchmark Limitations**: The model excels in math and coding but has room for improvement in other areas, such as common sense reasoning and nuanced language understanding.
const QwenQwQ32BPreview = "qwen/qwq-32b-preview"

// GoogleLearnLM15ProExperimentalFree is the ID for model Google: LearnLM 1.5 Pro Experimental (free)
//
// An experimental version of [Gemini 1.5 Pro](/google/gemini-pro-1.5) from Google.
const GoogleLearnLM15ProExperimentalFree = "google/learnlm-1.5-pro-experimental:free"

// EVAQwen2572B is the ID for model EVA Qwen2.5 72B
//
// EVA Qwen2.5 72B is a roleplay and storywriting specialist model. It's a full-parameter finetune of Qwen2.5-72B on mixture of synthetic and natural data.
//
// It uses Celeste 70B 0.1 data mixture, greatly expanding it to improve versatility, creativity and "flavor" of the resulting model.
const EVAQwen2572B = "eva-unit-01/eva-qwen-2.5-72b"

// OpenAIGPT4o20241120 is the ID for model OpenAI: GPT-4o (2024-11-20)
//
// The 2024-11-20 version of GPT-4o offers a leveled-up creative writing ability with more natural, engaging, and tailored writing to improve relevance & readability. It’s also better at working with uploaded files, providing deeper insights & more thorough responses.
//
// GPT-4o ("o" for "omni") is OpenAI's latest AI model, supporting both text and image inputs with text outputs. It maintains the intelligence level of [GPT-4 Turbo](/models/openai/gpt-4-turbo) while being twice as fast and 50% more cost-effective. GPT-4o also offers improved performance in processing non-English languages and enhanced visual capabilities.
const OpenAIGPT4o20241120 = "openai/gpt-4o-2024-11-20"

// MistralLarge2411 is the ID for model Mistral Large 2411
//
// Mistral Large 2 2411 is an update of [Mistral Large 2](/mistralai/mistral-large) released together with [Pixtral Large 2411](/mistralai/pixtral-large-2411)
//
// It provides a significant upgrade on the previous [Mistral Large 24.07](/mistralai/mistral-large-2407), with notable improvements in long context understanding, a new system prompt, and more accurate function calling.
const MistralLarge2411 = "mistralai/mistral-large-2411"

// MistralLarge2407 is the ID for model Mistral Large 2407
//
// This is Mistral AI's flagship model, Mistral Large 2 (version mistral-large-2407). It's a proprietary weights-available model and excels at reasoning, code, JSON, chat, and more. Read the launch announcement [here](https://mistral.ai/news/mistral-large-2407/).
//
// It supports dozens of languages including French, German, Spanish, Italian, Portuguese, Arabic, Hindi, Russian, Chinese, Japanese, and Korean, along with 80+ coding languages including Python, Java, C, C++, JavaScript, and Bash. Its long context window allows precise information recall from large documents.
const MistralLarge2407 = "mistralai/mistral-large-2407"

// MistralPixtralLarge2411 is the ID for model Mistral: Pixtral Large 2411
//
// Pixtral Large is a 124B parameter, open-weight, multimodal model built on top of [Mistral Large 2](/mistralai/mistral-large-2411). The model is able to understand documents, charts and natural images.
//
// The model is available under the Mistral Research License (MRL) for research and educational use, and the Mistral Commercial License for experimentation, testing, and production for commercial purposes.
const MistralPixtralLarge2411 = "mistralai/pixtral-large-2411"

// XAIGrokVisionBeta is the ID for model xAI: Grok Vision Beta
//
// Grok Vision Beta is xAI's experimental language model with vision capability.
const XAIGrokVisionBeta = "x-ai/grok-vision-beta"

// InfermaticMistralNemoInferor12B is the ID for model Infermatic: Mistral Nemo Inferor 12B
//
// Inferor 12B is a merge of top roleplay models, expert on immersive narratives and storytelling.
//
// This model was merged using the [Model Stock](https://arxiv.org/abs/2403.19522) merge method using [anthracite-org/magnum-v4-12b](https://openrouter.ai/anthracite-org/magnum-v4-72b) as a base.
const InfermaticMistralNemoInferor12B = "infermatic/mn-inferor-12b"

// Qwen25Coder32BInstruct is the ID for model Qwen2.5 Coder 32B Instruct
//
// Qwen2.5-Coder is the latest series of Code-Specific Qwen large language models (formerly known as CodeQwen). Qwen2.5-Coder brings the following improvements upon CodeQwen1.5:
//
// - Significantly improvements in **code generation**, **code reasoning** and **code fixing**.
// - A more comprehensive foundation for real-world applications such as **Code Agents**. Not only enhancing coding capabilities but also maintaining its strengths in mathematics and general competencies.
//
// To read more about its evaluation results, check out [Qwen 2.5 Coder's blog](https://qwenlm.github.io/blog/qwen2.5-coder-family/).
const Qwen25Coder32BInstruct = "qwen/qwen-2.5-coder-32b-instruct"

// SorcererLM8x22B is the ID for model SorcererLM 8x22B
//
// SorcererLM is an advanced RP and storytelling model, built as a Low-rank 16-bit LoRA fine-tuned on [WizardLM-2 8x22B](/microsoft/wizardlm-2-8x22b).
//
// - Advanced reasoning and emotional intelligence for engaging and immersive interactions
// - Vivid writing capabilities enriched with spatial and contextual awareness
// - Enhanced narrative depth, promoting creative and dynamic storytelling
const SorcererLM8x22B = "raifle/sorcererlm-8x22b"

// EVAQwen2532B is the ID for model EVA Qwen2.5 32B
//
// EVA Qwen2.5 32B is a roleplaying/storywriting specialist model. It's a full-parameter finetune of Qwen2.5-32B on mixture of synthetic and natural data.
//
// It uses Celeste 70B 0.1 data mixture, greatly expanding it to improve versatility, creativity and "flavor" of the resulting model.
const EVAQwen2532B = "eva-unit-01/eva-qwen-2.5-32b"

// Unslopnemo12b is the ID for model Unslopnemo 12b
//
// UnslopNemo v4.1 is the latest addition from the creator of Rocinante, designed for adventure writing and role-play scenarios.
const Unslopnemo12b = "thedrummer/unslopnemo-12b"

// AnthropicClaude35Haiku20241022SelfModerated is the ID for model Anthropic: Claude 3.5 Haiku (2024-10-22) (self-moderated)
//
// Claude 3.5 Haiku features enhancements across all skill sets including coding, tool use, and reasoning. As the fastest model in the Anthropic lineup, it offers rapid response times suitable for applications that require high interactivity and low latency, such as user-facing chatbots and on-the-fly code completions. It also excels in specialized tasks like data extraction and real-time content moderation, making it a versatile tool for a broad range of industries.
//
// It does not support image inputs.
//
// See the launch announcement and benchmark results [here](https://www.anthropic.com/news/3-5-models-and-computer-use)
const AnthropicClaude35Haiku20241022SelfModerated = "anthropic/claude-3.5-haiku-20241022:beta"

// AnthropicClaude35Haiku20241022 is the ID for model Anthropic: Claude 3.5 Haiku (2024-10-22)
//
// Claude 3.5 Haiku features enhancements across all skill sets including coding, tool use, and reasoning. As the fastest model in the Anthropic lineup, it offers rapid response times suitable for applications that require high interactivity and low latency, such as user-facing chatbots and on-the-fly code completions. It also excels in specialized tasks like data extraction and real-time content moderation, making it a versatile tool for a broad range of industries.
//
// It does not support image inputs.
//
// See the launch announcement and benchmark results [here](https://www.anthropic.com/news/3-5-models-and-computer-use)
const AnthropicClaude35Haiku20241022 = "anthropic/claude-3.5-haiku-20241022"

// AnthropicClaude35HaikuSelfModerated is the ID for model Anthropic: Claude 3.5 Haiku (self-moderated)
//
// Claude 3.5 Haiku features offers enhanced capabilities in speed, coding accuracy, and tool use. Engineered to excel in real-time applications, it delivers quick response times that are essential for dynamic tasks such as chat interactions and immediate coding suggestions.
//
// This makes it highly suitable for environments that demand both speed and precision, such as software development, customer service bots, and data management systems.
//
// This model is currently pointing to [Claude 3.5 Haiku (2024-10-22)](/anthropic/claude-3-5-haiku-20241022).
const AnthropicClaude35HaikuSelfModerated = "anthropic/claude-3.5-haiku:beta"

// AnthropicClaude35Haiku is the ID for model Anthropic: Claude 3.5 Haiku
//
// Claude 3.5 Haiku features offers enhanced capabilities in speed, coding accuracy, and tool use. Engineered to excel in real-time applications, it delivers quick response times that are essential for dynamic tasks such as chat interactions and immediate coding suggestions.
//
// This makes it highly suitable for environments that demand both speed and precision, such as software development, customer service bots, and data management systems.
//
// This model is currently pointing to [Claude 3.5 Haiku (2024-10-22)](/anthropic/claude-3-5-haiku-20241022).
const AnthropicClaude35Haiku = "anthropic/claude-3.5-haiku"

// NeverSleepLumimaidV0270B is the ID for model NeverSleep: Lumimaid v0.2 70B
//
// Lumimaid v0.2 70B is a finetune of [Llama 3.1 70B](/meta-llama/llama-3.1-70b-instruct) with a "HUGE step up dataset wise" compared to Lumimaid v0.1. Sloppy chats output were purged.
//
// Usage of this model is subject to [Meta's Acceptable Use Policy](https://llama.meta.com/llama3/use-policy/).
const NeverSleepLumimaidV0270B = "neversleep/llama-3.1-lumimaid-70b"

// MagnumV472B is the ID for model Magnum v4 72B
//
// This is a series of models designed to replicate the prose quality of the Claude 3 models, specifically Sonnet(https://openrouter.ai/anthropic/claude-3.5-sonnet) and Opus(https://openrouter.ai/anthropic/claude-3-opus).
//
// The model is fine-tuned on top of [Qwen2.5 72B](https://openrouter.ai/qwen/qwen-2.5-72b-instruct).
const MagnumV472B = "anthracite-org/magnum-v4-72b"

// AnthropicClaude35Sonnet is the ID for model Anthropic: Claude 3.5 Sonnet
//
// New Claude 3.5 Sonnet delivers better-than-Opus capabilities, faster-than-Sonnet speeds, at the same Sonnet prices. Sonnet is particularly good at:
//
// - Coding: Scores ~49% on SWE-Bench Verified, higher than the last best score, and without any fancy prompt scaffolding
// - Data science: Augments human data science expertise; navigates unstructured data while using multiple tools for insights
// - Visual processing: excelling at interpreting charts, graphs, and images, accurately transcribing text to derive insights beyond just the text alone
// - Agentic tasks: exceptional tool use, making it great at agentic tasks (i.e. complex, multi-step problem solving tasks that require engaging with other systems)
//
// #multimodal
const AnthropicClaude35Sonnet = "anthropic/claude-3.5-sonnet"

// AnthropicClaude35SonnetSelfModerated is the ID for model Anthropic: Claude 3.5 Sonnet (self-moderated)
//
// New Claude 3.5 Sonnet delivers better-than-Opus capabilities, faster-than-Sonnet speeds, at the same Sonnet prices. Sonnet is particularly good at:
//
// - Coding: Scores ~49% on SWE-Bench Verified, higher than the last best score, and without any fancy prompt scaffolding
// - Data science: Augments human data science expertise; navigates unstructured data while using multiple tools for insights
// - Visual processing: excelling at interpreting charts, graphs, and images, accurately transcribing text to derive insights beyond just the text alone
// - Agentic tasks: exceptional tool use, making it great at agentic tasks (i.e. complex, multi-step problem solving tasks that require engaging with other systems)
//
// #multimodal
const AnthropicClaude35SonnetSelfModerated = "anthropic/claude-3.5-sonnet:beta"

// XAIGrokBeta is the ID for model xAI: Grok Beta
//
// Grok Beta is xAI's experimental language model with state-of-the-art reasoning capabilities, best for complex and multi-step use cases.
//
// It is the successor of [Grok 2](https://x.ai/blog/grok-2) with enhanced context length.
const XAIGrokBeta = "x-ai/grok-beta"

// MistralMinistral8B is the ID for model Mistral: Ministral 8B
//
// Ministral 8B is an 8B parameter model featuring a unique interleaved sliding-window attention pattern for faster, memory-efficient inference. Designed for edge use cases, it supports up to 128k context length and excels in knowledge and reasoning tasks. It outperforms peers in the sub-10B category, making it perfect for low-latency, privacy-first applications.
const MistralMinistral8B = "mistralai/ministral-8b"

// MistralMinistral3B is the ID for model Mistral: Ministral 3B
//
// Ministral 3B is a 3B parameter model optimized for on-device and edge computing. It excels in knowledge, commonsense reasoning, and function-calling, outperforming larger models like Mistral 7B on most benchmarks. Supporting up to 128k context length, it’s ideal for orchestrating agentic workflows and specialist tasks with efficient inference.
const MistralMinistral3B = "mistralai/ministral-3b"

// Qwen257BInstruct is the ID for model Qwen2.5 7B Instruct
//
// Qwen2.5 7B is the latest series of Qwen large language models. Qwen2.5 brings the following improvements upon Qwen2:
//
// - Significantly more knowledge and has greatly improved capabilities in coding and mathematics, thanks to our specialized expert models in these domains.
//
// - Significant improvements in instruction following, generating long texts (over 8K tokens), understanding structured data (e.g, tables), and generating structured outputs especially JSON. More resilient to the diversity of system prompts, enhancing role-play implementation and condition-setting for chatbots.
//
// - Long-context Support up to 128K tokens and can generate up to 8K tokens.
//
// - Multilingual support for over 29 languages, including Chinese, English, French, Spanish, Portuguese, German, Italian, Russian, Japanese, Korean, Vietnamese, Thai, Arabic, and more.
//
// Usage of this model is subject to [Tongyi Qianwen LICENSE AGREEMENT](https://huggingface.co/Qwen/Qwen1.5-110B-Chat/blob/main/LICENSE).
const Qwen257BInstruct = "qwen/qwen-2.5-7b-instruct"

// NVIDIALlama31Nemotron70BInstructFree is the ID for model NVIDIA: Llama 3.1 Nemotron 70B Instruct (free)
//
// NVIDIA's Llama 3.1 Nemotron 70B is a language model designed for generating precise and useful responses. Leveraging [Llama 3.1 70B](/models/meta-llama/llama-3.1-70b-instruct) architecture and Reinforcement Learning from Human Feedback (RLHF), it excels in automatic alignment benchmarks. This model is tailored for applications requiring high accuracy in helpfulness and response generation, suitable for diverse user queries across multiple domains.
//
// Usage of this model is subject to [Meta's Acceptable Use Policy](https://www.llama.com/llama3/use-policy/).
const NVIDIALlama31Nemotron70BInstructFree = "nvidia/llama-3.1-nemotron-70b-instruct:free"

// NVIDIALlama31Nemotron70BInstruct is the ID for model NVIDIA: Llama 3.1 Nemotron 70B Instruct
//
// NVIDIA's Llama 3.1 Nemotron 70B is a language model designed for generating precise and useful responses. Leveraging [Llama 3.1 70B](/models/meta-llama/llama-3.1-70b-instruct) architecture and Reinforcement Learning from Human Feedback (RLHF), it excels in automatic alignment benchmarks. This model is tailored for applications requiring high accuracy in helpfulness and response generation, suitable for diverse user queries across multiple domains.
//
// Usage of this model is subject to [Meta's Acceptable Use Policy](https://www.llama.com/llama3/use-policy/).
const NVIDIALlama31Nemotron70BInstruct = "nvidia/llama-3.1-nemotron-70b-instruct"

// InflectionInflection3Pi is the ID for model Inflection: Inflection 3 Pi
//
// Inflection 3 Pi powers Inflection's [Pi](https://pi.ai) chatbot, including backstory, emotional intelligence, productivity, and safety. It has access to recent news, and excels in scenarios like customer support and roleplay.
//
// Pi has been trained to mirror your tone and style, if you use more emojis, so will Pi! Try experimenting with various prompts and conversation styles.
const InflectionInflection3Pi = "inflection/inflection-3-pi"

// InflectionInflection3Productivity is the ID for model Inflection: Inflection 3 Productivity
//
// Inflection 3 Productivity is optimized for following instructions. It is better for tasks requiring JSON output or precise adherence to provided guidelines. It has access to recent news.
//
// For emotional intelligence similar to Pi, see [Inflect 3 Pi](/inflection/inflection-3-pi)
//
// See [Inflection's announcement](https://inflection.ai/blog/enterprise) for more details.
const InflectionInflection3Productivity = "inflection/inflection-3-productivity"

// GoogleGeminiFlash158B is the ID for model Google: Gemini Flash 1.5 8B
//
// Gemini Flash 1.5 8B is optimized for speed and efficiency, offering enhanced performance in small prompt tasks like chat, transcription, and translation. With reduced latency, it is highly effective for real-time and large-scale operations. This model focuses on cost-effective solutions while maintaining high-quality results.
//
// [Click here to learn more about this model](https://developers.googleblog.com/en/gemini-15-flash-8b-is-now-generally-available-for-use/).
//
// Usage of Gemini is subject to Google's [Gemini Terms of Use](https://ai.google.dev/terms).
const GoogleGeminiFlash158B = "google/gemini-flash-1.5-8b"

// MagnumV272B is the ID for model Magnum v2 72B
//
// From the maker of [Goliath](https://openrouter.ai/models/alpindale/goliath-120b), Magnum 72B is the seventh in a family of models designed to achieve the prose quality of the Claude 3 models, notably Opus & Sonnet.
//
// The model is based on [Qwen2 72B](https://openrouter.ai/models/qwen/qwen-2-72b-instruct) and trained with 55 million tokens of highly curated roleplay (RP) data.
const MagnumV272B = "anthracite-org/magnum-v2-72b"

// LiquidLFM40BMoE is the ID for model Liquid: LFM 40B MoE
//
// Liquid's 40.3B Mixture of Experts (MoE) model. Liquid Foundation Models (LFMs) are large neural networks built with computational units rooted in dynamic systems.
//
// LFMs are general-purpose AI models that can be used to model any kind of sequential data, including video, audio, text, time series, and signals.
//
// See the [launch announcement](https://www.liquid.ai/liquid-foundation-models) for benchmarks and more info.
const LiquidLFM40BMoE = "liquid/lfm-40b"

// Rocinante12B is the ID for model Rocinante 12B
//
// Rocinante 12B is designed for engaging storytelling and rich prose.
//
// Early testers have reported:
// - Expanded vocabulary with unique and expressive word choices
// - Enhanced creativity for vivid narratives
// - Adventure-filled and captivating stories
const Rocinante12B = "thedrummer/rocinante-12b"

// MetaLlama323BInstruct is the ID for model Meta: Llama 3.2 3B Instruct
//
// Llama 3.2 3B is a 3-billion-parameter multilingual large language model, optimized for advanced natural language processing tasks like dialogue generation, reasoning, and summarization. Designed with the latest transformer architecture, it supports eight languages, including English, Spanish, and Hindi, and is adaptable for additional languages.
//
// Trained on 9 trillion tokens, the Llama 3.2 3B model excels in instruction-following, complex reasoning, and tool use. Its balanced performance makes it ideal for applications needing accuracy and efficiency in text generation across multilingual settings.
//
// Click here for the [original model card](https://github.com/meta-llama/llama-models/blob/main/models/llama3_2/CARD.md).
//
// Usage of this model is subject to [Meta's Acceptable Use Policy](https://www.llama.com/llama3/use-policy/).
const MetaLlama323BInstruct = "meta-llama/llama-3.2-3b-instruct"

// MetaLlama321BInstruct is the ID for model Meta: Llama 3.2 1B Instruct
//
// Llama 3.2 1B is a 1-billion-parameter language model focused on efficiently performing natural language tasks, such as summarization, dialogue, and multilingual text analysis. Its smaller size allows it to operate efficiently in low-resource environments while maintaining strong task performance.
//
// Supporting eight core languages and fine-tunable for more, Llama 1.3B is ideal for businesses or developers seeking lightweight yet powerful AI solutions that can operate in diverse multilingual settings without the high computational demand of larger models.
//
// Click here for the [original model card](https://github.com/meta-llama/llama-models/blob/main/models/llama3_2/CARD.md).
//
// Usage of this model is subject to [Meta's Acceptable Use Policy](https://www.llama.com/llama3/use-policy/).
const MetaLlama321BInstruct = "meta-llama/llama-3.2-1b-instruct"

// MetaLlama3290BVisionInstruct is the ID for model Meta: Llama 3.2 90B Vision Instruct
//
// The Llama 90B Vision model is a top-tier, 90-billion-parameter multimodal model designed for the most challenging visual reasoning and language tasks. It offers unparalleled accuracy in image captioning, visual question answering, and advanced image-text comprehension. Pre-trained on vast multimodal datasets and fine-tuned with human feedback, the Llama 90B Vision is engineered to handle the most demanding image-based AI tasks.
//
// This model is perfect for industries requiring cutting-edge multimodal AI capabilities, particularly those dealing with complex, real-time visual and textual analysis.
//
// Click here for the [original model card](https://github.com/meta-llama/llama-models/blob/main/models/llama3_2/CARD_VISION.md).
//
// Usage of this model is subject to [Meta's Acceptable Use Policy](https://www.llama.com/llama3/use-policy/).
const MetaLlama3290BVisionInstruct = "meta-llama/llama-3.2-90b-vision-instruct"

// MetaLlama3211BVisionInstructFree is the ID for model Meta: Llama 3.2 11B Vision Instruct (free)
//
// Llama 3.2 11B Vision is a multimodal model with 11 billion parameters, designed to handle tasks combining visual and textual data. It excels in tasks such as image captioning and visual question answering, bridging the gap between language generation and visual reasoning. Pre-trained on a massive dataset of image-text pairs, it performs well in complex, high-accuracy image analysis.
//
// Its ability to integrate visual understanding with language processing makes it an ideal solution for industries requiring comprehensive visual-linguistic AI applications, such as content creation, AI-driven customer service, and research.
//
// Click here for the [original model card](https://github.com/meta-llama/llama-models/blob/main/models/llama3_2/CARD_VISION.md).
//
// Usage of this model is subject to [Meta's Acceptable Use Policy](https://www.llama.com/llama3/use-policy/).
const MetaLlama3211BVisionInstructFree = "meta-llama/llama-3.2-11b-vision-instruct:free"

// MetaLlama3211BVisionInstruct is the ID for model Meta: Llama 3.2 11B Vision Instruct
//
// Llama 3.2 11B Vision is a multimodal model with 11 billion parameters, designed to handle tasks combining visual and textual data. It excels in tasks such as image captioning and visual question answering, bridging the gap between language generation and visual reasoning. Pre-trained on a massive dataset of image-text pairs, it performs well in complex, high-accuracy image analysis.
//
// Its ability to integrate visual understanding with language processing makes it an ideal solution for industries requiring comprehensive visual-linguistic AI applications, such as content creation, AI-driven customer service, and research.
//
// Click here for the [original model card](https://github.com/meta-llama/llama-models/blob/main/models/llama3_2/CARD_VISION.md).
//
// Usage of this model is subject to [Meta's Acceptable Use Policy](https://www.llama.com/llama3/use-policy/).
const MetaLlama3211BVisionInstruct = "meta-llama/llama-3.2-11b-vision-instruct"

// Qwen2572BInstruct is the ID for model Qwen2.5 72B Instruct
//
// Qwen2.5 72B is the latest series of Qwen large language models. Qwen2.5 brings the following improvements upon Qwen2:
//
// - Significantly more knowledge and has greatly improved capabilities in coding and mathematics, thanks to our specialized expert models in these domains.
//
// - Significant improvements in instruction following, generating long texts (over 8K tokens), understanding structured data (e.g, tables), and generating structured outputs especially JSON. More resilient to the diversity of system prompts, enhancing role-play implementation and condition-setting for chatbots.
//
// - Long-context Support up to 128K tokens and can generate up to 8K tokens.
//
// - Multilingual support for over 29 languages, including Chinese, English, French, Spanish, Portuguese, German, Italian, Russian, Japanese, Korean, Vietnamese, Thai, Arabic, and more.
//
// Usage of this model is subject to [Tongyi Qianwen LICENSE AGREEMENT](https://huggingface.co/Qwen/Qwen1.5-110B-Chat/blob/main/LICENSE).
const Qwen2572BInstruct = "qwen/qwen-2.5-72b-instruct"

// Qwen2VL72BInstruct is the ID for model Qwen2-VL 72B Instruct
//
// Qwen2 VL 72B is a multimodal LLM from the Qwen Team with the following key enhancements:
//
// - SoTA understanding of images of various resolution & ratio: Qwen2-VL achieves state-of-the-art performance on visual understanding benchmarks, including MathVista, DocVQA, RealWorldQA, MTVQA, etc.
//
// - Understanding videos of 20min+: Qwen2-VL can understand videos over 20 minutes for high-quality video-based question answering, dialog, content creation, etc.
//
// - Agent that can operate your mobiles, robots, etc.: with the abilities of complex reasoning and decision making, Qwen2-VL can be integrated with devices like mobile phones, robots, etc., for automatic operation based on visual environment and text instructions.
//
// - Multilingual Support: to serve global users, besides English and Chinese, Qwen2-VL now supports the understanding of texts in different languages inside images, including most European languages, Japanese, Korean, Arabic, Vietnamese, etc.
//
// For more details, see this [blog post](https://qwenlm.github.io/blog/qwen2-vl/) and [GitHub repo](https://github.com/QwenLM/Qwen2-VL).
//
// Usage of this model is subject to [Tongyi Qianwen LICENSE AGREEMENT](https://huggingface.co/Qwen/Qwen1.5-110B-Chat/blob/main/LICENSE).
const Qwen2VL72BInstruct = "qwen/qwen-2-vl-72b-instruct"

// NeverSleepLumimaidV028B is the ID for model NeverSleep: Lumimaid v0.2 8B
//
// Lumimaid v0.2 8B is a finetune of [Llama 3.1 8B](/models/meta-llama/llama-3.1-8b-instruct) with a "HUGE step up dataset wise" compared to Lumimaid v0.1. Sloppy chats output were purged.
//
// Usage of this model is subject to [Meta's Acceptable Use Policy](https://llama.meta.com/llama3/use-policy/).
const NeverSleepLumimaidV028B = "neversleep/llama-3.1-lumimaid-8b"

// OpenAIO1Mini20240912 is the ID for model OpenAI: o1-mini (2024-09-12)
//
// The latest and strongest model family from OpenAI, o1 is designed to spend more time thinking before responding.
//
// The o1 models are optimized for math, science, programming, and other STEM-related tasks. They consistently exhibit PhD-level accuracy on benchmarks in physics, chemistry, and biology. Learn more in the [launch announcement](https://openai.com/o1).
//
// Note: This model is currently experimental and not suitable for production use-cases, and may be heavily rate-limited.
const OpenAIO1Mini20240912 = "openai/o1-mini-2024-09-12"

// OpenAIO1Preview is the ID for model OpenAI: o1-preview
//
// The latest and strongest model family from OpenAI, o1 is designed to spend more time thinking before responding.
//
// The o1 models are optimized for math, science, programming, and other STEM-related tasks. They consistently exhibit PhD-level accuracy on benchmarks in physics, chemistry, and biology. Learn more in the [launch announcement](https://openai.com/o1).
//
// Note: This model is currently experimental and not suitable for production use-cases, and may be heavily rate-limited.
const OpenAIO1Preview = "openai/o1-preview"

// OpenAIO1Preview20240912 is the ID for model OpenAI: o1-preview (2024-09-12)
//
// The latest and strongest model family from OpenAI, o1 is designed to spend more time thinking before responding.
//
// The o1 models are optimized for math, science, programming, and other STEM-related tasks. They consistently exhibit PhD-level accuracy on benchmarks in physics, chemistry, and biology. Learn more in the [launch announcement](https://openai.com/o1).
//
// Note: This model is currently experimental and not suitable for production use-cases, and may be heavily rate-limited.
const OpenAIO1Preview20240912 = "openai/o1-preview-2024-09-12"

// OpenAIO1Mini is the ID for model OpenAI: o1-mini
//
// The latest and strongest model family from OpenAI, o1 is designed to spend more time thinking before responding.
//
// The o1 models are optimized for math, science, programming, and other STEM-related tasks. They consistently exhibit PhD-level accuracy on benchmarks in physics, chemistry, and biology. Learn more in the [launch announcement](https://openai.com/o1).
//
// Note: This model is currently experimental and not suitable for production use-cases, and may be heavily rate-limited.
const OpenAIO1Mini = "openai/o1-mini"

// MistralPixtral12B is the ID for model Mistral: Pixtral 12B
//
// The first multi-modal, text+image-to-text model from Mistral AI. Its weights were launched via torrent: https://x.com/mistralai/status/1833758285167722836.
const MistralPixtral12B = "mistralai/pixtral-12b"

// CohereCommandR082024 is the ID for model Cohere: Command R (08-2024)
//
// command-r-08-2024 is an update of the [Command R](/models/cohere/command-r) with improved performance for multilingual retrieval-augmented generation (RAG) and tool use. More broadly, it is better at math, code and reasoning and is competitive with the previous version of the larger Command R+ model.
//
// Read the launch post [here](https://docs.cohere.com/changelog/command-gets-refreshed).
//
// Use of this model is subject to Cohere's [Usage Policy](https://docs.cohere.com/docs/usage-policy) and [SaaS Agreement](https://cohere.com/saas-agreement).
const CohereCommandR082024 = "cohere/command-r-08-2024"

// CohereCommandR082024Plus is the ID for model Cohere: Command R+ (08-2024)
//
// command-r-plus-08-2024 is an update of the [Command R+](/models/cohere/command-r-plus) with roughly 50% higher throughput and 25% lower latencies as compared to the previous Command R+ version, while keeping the hardware footprint the same.
//
// Read the launch post [here](https://docs.cohere.com/changelog/command-gets-refreshed).
//
// Use of this model is subject to Cohere's [Usage Policy](https://docs.cohere.com/docs/usage-policy) and [SaaS Agreement](https://cohere.com/saas-agreement).
const CohereCommandR082024Plus = "cohere/command-r-plus-08-2024"

// Qwen2VL7BInstruct is the ID for model Qwen2-VL 7B Instruct
//
// Qwen2 VL 7B is a multimodal LLM from the Qwen Team with the following key enhancements:
//
// - SoTA understanding of images of various resolution & ratio: Qwen2-VL achieves state-of-the-art performance on visual understanding benchmarks, including MathVista, DocVQA, RealWorldQA, MTVQA, etc.
//
// - Understanding videos of 20min+: Qwen2-VL can understand videos over 20 minutes for high-quality video-based question answering, dialog, content creation, etc.
//
// - Agent that can operate your mobiles, robots, etc.: with the abilities of complex reasoning and decision making, Qwen2-VL can be integrated with devices like mobile phones, robots, etc., for automatic operation based on visual environment and text instructions.
//
// - Multilingual Support: to serve global users, besides English and Chinese, Qwen2-VL now supports the understanding of texts in different languages inside images, including most European languages, Japanese, Korean, Arabic, Vietnamese, etc.
//
// For more details, see this [blog post](https://qwenlm.github.io/blog/qwen2-vl/) and [GitHub repo](https://github.com/QwenLM/Qwen2-VL).
//
// Usage of this model is subject to [Tongyi Qianwen LICENSE AGREEMENT](https://huggingface.co/Qwen/Qwen1.5-110B-Chat/blob/main/LICENSE).
const Qwen2VL7BInstruct = "qwen/qwen-2-vl-7b-instruct"

// Sao10KLlama31Euryale70BV22 is the ID for model Sao10K: Llama 3.1 Euryale 70B v2.2
//
// Euryale L3.1 70B v2.2 is a model focused on creative roleplay from [Sao10k](https://ko-fi.com/sao10k). It is the successor of [Euryale L3 70B v2.1](/models/sao10k/l3-euryale-70b).
const Sao10KLlama31Euryale70BV22 = "sao10k/l3.1-euryale-70b"

// GoogleGeminiFlash158BExperimental is the ID for model Google: Gemini Flash 1.5 8B Experimental
//
// Gemini Flash 1.5 8B Experimental is an experimental, 8B parameter version of the [Gemini Flash 1.5](/models/google/gemini-flash-1.5) model.
//
// Usage of Gemini is subject to Google's [Gemini Terms of Use](https://ai.google.dev/terms).
//
// #multimodal
//
// Note: This model is currently experimental and not suitable for production use-cases, and may be heavily rate-limited.
const GoogleGeminiFlash158BExperimental = "google/gemini-flash-1.5-8b-exp"

// AI21Jamba15Large is the ID for model AI21: Jamba 1.5 Large
//
// Jamba 1.5 Large is part of AI21's new family of open models, offering superior speed, efficiency, and quality.
//
// It features a 256K effective context window, the longest among open models, enabling improved performance on tasks like document summarization and analysis.
//
// Built on a novel SSM-Transformer architecture, it outperforms larger models like Llama 3.1 70B on benchmarks while maintaining resource efficiency.
//
// Read their [announcement](https://www.ai21.com/blog/announcing-jamba-model-family) to learn more.
const AI21Jamba15Large = "ai21/jamba-1-5-large"

// AI21Jamba15Mini is the ID for model AI21: Jamba 1.5 Mini
//
// Jamba 1.5 Mini is the world's first production-grade Mamba-based model, combining SSM and Transformer architectures for a 256K context window and high efficiency.
//
// It works with 9 languages and can handle various writing and analysis tasks as well as or better than similar small models.
//
// This model uses less computer memory and works faster with longer texts than previous designs.
//
// Read their [announcement](https://www.ai21.com/blog/announcing-jamba-model-family) to learn more.
const AI21Jamba15Mini = "ai21/jamba-1-5-mini"

// MicrosoftPhi35Mini128KInstruct is the ID for model Microsoft: Phi-3.5 Mini 128K Instruct
//
// Phi-3.5 models are lightweight, state-of-the-art open models. These models were trained with Phi-3 datasets that include both synthetic data and the filtered, publicly available websites data, with a focus on high quality and reasoning-dense properties. Phi-3.5 Mini uses 3.8B parameters, and is a dense decoder-only transformer model using the same tokenizer as [Phi-3 Mini](/models/microsoft/phi-3-mini-128k-instruct).
//
// The models underwent a rigorous enhancement process, incorporating both supervised fine-tuning, proximal policy optimization, and direct preference optimization to ensure precise instruction adherence and robust safety measures. When assessed against benchmarks that test common sense, language understanding, math, code, long context and logical reasoning, Phi-3.5 models showcased robust and state-of-the-art performance among models with less than 13 billion parameters.
const MicrosoftPhi35Mini128KInstruct = "microsoft/phi-3.5-mini-128k-instruct"

// NousHermes370BInstruct is the ID for model Nous: Hermes 3 70B Instruct
//
// Hermes 3 is a generalist language model with many improvements over [Hermes 2](/models/nousresearch/nous-hermes-2-mistral-7b-dpo), including advanced agentic capabilities, much better roleplaying, reasoning, multi-turn conversation, long context coherence, and improvements across the board.
//
// Hermes 3 70B is a competitive, if not superior finetune of the [Llama-3.1 70B foundation model](/models/meta-llama/llama-3.1-70b-instruct), focused on aligning LLMs to the user, with powerful steering capabilities and control given to the end user.
//
// The Hermes 3 series builds and expands on the Hermes 2 set of capabilities, including more powerful and reliable function calling and structured output capabilities, generalist assistant capabilities, and improved code generation skills.
const NousHermes370BInstruct = "nousresearch/hermes-3-llama-3.1-70b"

// NousHermes3405BInstruct is the ID for model Nous: Hermes 3 405B Instruct
//
// Hermes 3 is a generalist language model with many improvements over Hermes 2, including advanced agentic capabilities, much better roleplaying, reasoning, multi-turn conversation, long context coherence, and improvements across the board.
//
// Hermes 3 405B is a frontier-level, full-parameter finetune of the Llama-3.1 405B foundation model, focused on aligning LLMs to the user, with powerful steering capabilities and control given to the end user.
//
// The Hermes 3 series builds and expands on the Hermes 2 set of capabilities, including more powerful and reliable function calling and structured output capabilities, generalist assistant capabilities, and improved code generation skills.
//
// Hermes 3 is competitive, if not superior, to Llama-3.1 Instruct models at general capabilities, with varying strengths and weaknesses attributable between the two.
const NousHermes3405BInstruct = "nousresearch/hermes-3-llama-3.1-405b"

// PerplexityLlama31Sonar405BOnline is the ID for model Perplexity: Llama 3.1 Sonar 405B Online
//
// Llama 3.1 Sonar is Perplexity's latest model family. It surpasses their earlier Sonar models in cost-efficiency, speed, and performance. The model is built upon the Llama 3.1 405B and has internet access.
const PerplexityLlama31Sonar405BOnline = "perplexity/llama-3.1-sonar-huge-128k-online"

// OpenAIChatGPT4o is the ID for model OpenAI: ChatGPT-4o
//
// OpenAI ChatGPT 4o is continually updated by OpenAI to point to the current version of GPT-4o used by ChatGPT. It therefore differs slightly from the API version of [GPT-4o](/models/openai/gpt-4o) in that it has additional RLHF. It is intended for research and evaluation.
//
// OpenAI notes that this model is not suited for production use-cases as it may be removed or redirected to another model in the future.
const OpenAIChatGPT4o = "openai/chatgpt-4o-latest"

// Sao10KLlama38BLunaris is the ID for model Sao10K: Llama 3 8B Lunaris
//
// Lunaris 8B is a versatile generalist and roleplaying model based on Llama 3. It's a strategic merge of multiple models, designed to balance creativity with improved logic and general knowledge.
//
// Created by [Sao10k](https://huggingface.co/Sao10k), this model aims to offer an improved experience over Stheno v3.2, with enhanced creativity and logical reasoning.
//
// For best results, use with Llama 3 Instruct context template, temperature 1.4, and min_p 0.1.
const Sao10KLlama38BLunaris = "sao10k/l3-lunaris-8b"

// AetherwiingStarcannon12B is the ID for model Aetherwiing: Starcannon 12B
//
// Starcannon 12B v2 is a creative roleplay and story writing model, based on Mistral Nemo, using [nothingiisreal/mn-celeste-12b](/nothingiisreal/mn-celeste-12b) as a base, with [intervitens/mini-magnum-12b-v1.1](https://huggingface.co/intervitens/mini-magnum-12b-v1.1) merged in using the [TIES](https://arxiv.org/abs/2306.01708) method.
//
// Although more similar to Magnum overall, the model remains very creative, with a pleasant writing style. It is recommended for people wanting more variety than Magnum, and yet more verbose prose than Celeste.
const AetherwiingStarcannon12B = "aetherwiing/mn-starcannon-12b"

// OpenAIGPT4o20240806 is the ID for model OpenAI: GPT-4o (2024-08-06)
//
// The 2024-08-06 version of GPT-4o offers improved performance in structured outputs, with the ability to supply a JSON schema in the response_format. Read more [here](https://openai.com/index/introducing-structured-outputs-in-the-api/).
//
// GPT-4o ("o" for "omni") is OpenAI's latest AI model, supporting both text and image inputs with text outputs. It maintains the intelligence level of [GPT-4 Turbo](/models/openai/gpt-4-turbo) while being twice as fast and 50% more cost-effective. GPT-4o also offers improved performance in processing non-English languages and enhanced visual capabilities.
//
// For benchmarking against other models, it was briefly called ["im-also-a-good-gpt2-chatbot"](https://twitter.com/LiamFedus/status/1790064963966370209)
const OpenAIGPT4o20240806 = "openai/gpt-4o-2024-08-06"

// MetaLlama31405BBase is the ID for model Meta: Llama 3.1 405B (base)
//
// Meta's latest class of model (Llama 3.1) launched with a variety of sizes & flavors. This is the base 405B pre-trained version.
//
// It has demonstrated strong performance compared to leading closed-source models in human evaluations.
//
// To read more about the model release, [click here](https://ai.meta.com/blog/meta-llama-3/). Usage of this model is subject to [Meta's Acceptable Use Policy](https://llama.meta.com/llama3/use-policy/).
const MetaLlama31405BBase = "meta-llama/llama-3.1-405b"

// MistralNemo12BCeleste is the ID for model Mistral Nemo 12B Celeste
//
// A specialized story writing and roleplaying model based on Mistral's NeMo 12B Instruct. Fine-tuned on curated datasets including Reddit Writing Prompts and Opus Instruct 25K.
//
// This model excels at creative writing, offering improved NSFW capabilities, with smarter and more active narration. It demonstrates remarkable versatility in both SFW and NSFW scenarios, with strong Out of Character (OOC) steering capabilities, allowing fine-tuned control over narrative direction and character behavior.
//
// Check out the model's [HuggingFace page](https://huggingface.co/nothingiisreal/MN-12B-Celeste-V1.9) for details on what parameters and prompts work best!
const MistralNemo12BCeleste = "nothingiisreal/mn-celeste-12b"

// PerplexityLlama31Sonar8B is the ID for model Perplexity: Llama 3.1 Sonar 8B
//
// Llama 3.1 Sonar is Perplexity's latest model family. It surpasses their earlier Sonar models in cost-efficiency, speed, and performance.
//
// This is a normal offline LLM, but the [online version](/models/perplexity/llama-3.1-sonar-small-128k-online) of this model has Internet access.
const PerplexityLlama31Sonar8B = "perplexity/llama-3.1-sonar-small-128k-chat"

// PerplexityLlama31Sonar70B is the ID for model Perplexity: Llama 3.1 Sonar 70B
//
// Llama 3.1 Sonar is Perplexity's latest model family. It surpasses their earlier Sonar models in cost-efficiency, speed, and performance.
//
// This is a normal offline LLM, but the [online version](/models/perplexity/llama-3.1-sonar-large-128k-online) of this model has Internet access.
const PerplexityLlama31Sonar70B = "perplexity/llama-3.1-sonar-large-128k-chat"

// PerplexityLlama31Sonar70BOnline is the ID for model Perplexity: Llama 3.1 Sonar 70B Online
//
// Llama 3.1 Sonar is Perplexity's latest model family. It surpasses their earlier Sonar models in cost-efficiency, speed, and performance.
//
// This is the online version of the [offline chat model](/models/perplexity/llama-3.1-sonar-large-128k-chat). It is focused on delivering helpful, up-to-date, and factual responses. #online
const PerplexityLlama31Sonar70BOnline = "perplexity/llama-3.1-sonar-large-128k-online"

// PerplexityLlama31Sonar8BOnline is the ID for model Perplexity: Llama 3.1 Sonar 8B Online
//
// Llama 3.1 Sonar is Perplexity's latest model family. It surpasses their earlier Sonar models in cost-efficiency, speed, and performance.
//
// This is the online version of the [offline chat model](/models/perplexity/llama-3.1-sonar-small-128k-chat). It is focused on delivering helpful, up-to-date, and factual responses. #online
const PerplexityLlama31Sonar8BOnline = "perplexity/llama-3.1-sonar-small-128k-online"

// MetaLlama31405BInstruct is the ID for model Meta: Llama 3.1 405B Instruct
//
// The highly anticipated 400B class of Llama3 is here! Clocking in at 128k context with impressive eval scores, the Meta AI team continues to push the frontier of open-source LLMs.
//
// Meta's latest class of model (Llama 3.1) launched with a variety of sizes & flavors. This 405B instruct-tuned version is optimized for high quality dialogue usecases.
//
// It has demonstrated strong performance compared to leading closed-source models including GPT-4o and Claude 3.5 Sonnet in evaluations.
//
// To read more about the model release, [click here](https://ai.meta.com/blog/meta-llama-3-1/). Usage of this model is subject to [Meta's Acceptable Use Policy](https://llama.meta.com/llama3/use-policy/).
const MetaLlama31405BInstruct = "meta-llama/llama-3.1-405b-instruct"

// MetaLlama318BInstruct is the ID for model Meta: Llama 3.1 8B Instruct
//
// Meta's latest class of model (Llama 3.1) launched with a variety of sizes & flavors. This 8B instruct-tuned version is fast and efficient.
//
// It has demonstrated strong performance compared to leading closed-source models in human evaluations.
//
// To read more about the model release, [click here](https://ai.meta.com/blog/meta-llama-3-1/). Usage of this model is subject to [Meta's Acceptable Use Policy](https://llama.meta.com/llama3/use-policy/).
const MetaLlama318BInstruct = "meta-llama/llama-3.1-8b-instruct"

// MetaLlama3170BInstruct is the ID for model Meta: Llama 3.1 70B Instruct
//
// Meta's latest class of model (Llama 3.1) launched with a variety of sizes & flavors. This 70B instruct-tuned version is optimized for high quality dialogue usecases.
//
// It has demonstrated strong performance compared to leading closed-source models in human evaluations.
//
// To read more about the model release, [click here](https://ai.meta.com/blog/meta-llama-3-1/). Usage of this model is subject to [Meta's Acceptable Use Policy](https://llama.meta.com/llama3/use-policy/).
const MetaLlama3170BInstruct = "meta-llama/llama-3.1-70b-instruct"

// MistralMistralNemoFree is the ID for model Mistral: Mistral Nemo (free)
//
// A 12B parameter model with a 128k token context length built by Mistral in collaboration with NVIDIA.
//
// The model is multilingual, supporting English, French, German, Spanish, Italian, Portuguese, Chinese, Japanese, Korean, Arabic, and Hindi.
//
// It supports function calling and is released under the Apache 2.0 license.
const MistralMistralNemoFree = "mistralai/mistral-nemo:free"

// MistralMistralNemo is the ID for model Mistral: Mistral Nemo
//
// A 12B parameter model with a 128k token context length built by Mistral in collaboration with NVIDIA.
//
// The model is multilingual, supporting English, French, German, Spanish, Italian, Portuguese, Chinese, Japanese, Korean, Arabic, and Hindi.
//
// It supports function calling and is released under the Apache 2.0 license.
const MistralMistralNemo = "mistralai/mistral-nemo"

// MistralCodestralMamba is the ID for model Mistral: Codestral Mamba
//
// A 7.3B parameter Mamba-based model designed for code and reasoning tasks.
//
// - Linear time inference, allowing for theoretically infinite sequence lengths
// - 256k token context window
// - Optimized for quick responses, especially beneficial for code productivity
// - Performs comparably to state-of-the-art transformer models in code and reasoning tasks
// - Available under the Apache 2.0 license for free use, modification, and distribution
const MistralCodestralMamba = "mistralai/codestral-mamba"

// OpenAIGPT4oMini is the ID for model OpenAI: GPT-4o-mini
//
// GPT-4o mini is OpenAI's newest model after [GPT-4 Omni](/models/openai/gpt-4o), supporting both text and image inputs with text outputs.
//
// As their most advanced small model, it is many multiples more affordable than other recent frontier models, and more than 60% cheaper than [GPT-3.5 Turbo](/models/openai/gpt-3.5-turbo). It maintains SOTA intelligence, while being significantly more cost-effective.
//
// GPT-4o mini achieves an 82% score on MMLU and presently ranks higher than GPT-4 on chat preferences [common leaderboards](https://arena.lmsys.org/).
//
// Check out the [launch announcement](https://openai.com/index/gpt-4o-mini-advancing-cost-efficient-intelligence/) to learn more.
//
// #multimodal
const OpenAIGPT4oMini = "openai/gpt-4o-mini"

// OpenAIGPT4oMini20240718 is the ID for model OpenAI: GPT-4o-mini (2024-07-18)
//
// GPT-4o mini is OpenAI's newest model after [GPT-4 Omni](/models/openai/gpt-4o), supporting both text and image inputs with text outputs.
//
// As their most advanced small model, it is many multiples more affordable than other recent frontier models, and more than 60% cheaper than [GPT-3.5 Turbo](/models/openai/gpt-3.5-turbo). It maintains SOTA intelligence, while being significantly more cost-effective.
//
// GPT-4o mini achieves an 82% score on MMLU and presently ranks higher than GPT-4 on chat preferences [common leaderboards](https://arena.lmsys.org/).
//
// Check out the [launch announcement](https://openai.com/index/gpt-4o-mini-advancing-cost-efficient-intelligence/) to learn more.
//
// #multimodal
const OpenAIGPT4oMini20240718 = "openai/gpt-4o-mini-2024-07-18"

// Qwen27BInstructFree is the ID for model Qwen 2 7B Instruct (free)
//
// Qwen2 7B is a transformer-based model that excels in language understanding, multilingual capabilities, coding, mathematics, and reasoning.
//
// It features SwiGLU activation, attention QKV bias, and group query attention. It is pretrained on extensive data with supervised finetuning and direct preference optimization.
//
// For more details, see this [blog post](https://qwenlm.github.io/blog/qwen2/) and [GitHub repo](https://github.com/QwenLM/Qwen2).
//
// Usage of this model is subject to [Tongyi Qianwen LICENSE AGREEMENT](https://huggingface.co/Qwen/Qwen1.5-110B-Chat/blob/main/LICENSE).
const Qwen27BInstructFree = "qwen/qwen-2-7b-instruct:free"

// Qwen27BInstruct is the ID for model Qwen 2 7B Instruct
//
// Qwen2 7B is a transformer-based model that excels in language understanding, multilingual capabilities, coding, mathematics, and reasoning.
//
// It features SwiGLU activation, attention QKV bias, and group query attention. It is pretrained on extensive data with supervised finetuning and direct preference optimization.
//
// For more details, see this [blog post](https://qwenlm.github.io/blog/qwen2/) and [GitHub repo](https://github.com/QwenLM/Qwen2).
//
// Usage of this model is subject to [Tongyi Qianwen LICENSE AGREEMENT](https://huggingface.co/Qwen/Qwen1.5-110B-Chat/blob/main/LICENSE).
const Qwen27BInstruct = "qwen/qwen-2-7b-instruct"

// GoogleGemma227B is the ID for model Google: Gemma 2 27B
//
// Gemma 2 27B by Google is an open model built from the same research and technology used to create the [Gemini models](/models?q=gemini).
//
// Gemma models are well-suited for a variety of text generation tasks, including question answering, summarization, and reasoning.
//
// See the [launch announcement](https://blog.google/technology/developers/google-gemma-2/) for more details. Usage of Gemma is subject to Google's [Gemma Terms of Use](https://ai.google.dev/gemma/terms).
const GoogleGemma227B = "google/gemma-2-27b-it"

// Magnum72B is the ID for model Magnum 72B
//
// From the maker of [Goliath](https://openrouter.ai/models/alpindale/goliath-120b), Magnum 72B is the first in a new family of models designed to achieve the prose quality of the Claude 3 models, notably Opus & Sonnet.
//
// The model is based on [Qwen2 72B](https://openrouter.ai/models/qwen/qwen-2-72b-instruct) and trained with 55 million tokens of highly curated roleplay (RP) data.
const Magnum72B = "alpindale/magnum-72b"

// GoogleGemma29BFree is the ID for model Google: Gemma 2 9B (free)
//
// Gemma 2 9B by Google is an advanced, open-source language model that sets a new standard for efficiency and performance in its size class.
//
// Designed for a wide variety of tasks, it empowers developers and researchers to build innovative applications, while maintaining accessibility, safety, and cost-effectiveness.
//
// See the [launch announcement](https://blog.google/technology/developers/google-gemma-2/) for more details. Usage of Gemma is subject to Google's [Gemma Terms of Use](https://ai.google.dev/gemma/terms).
const GoogleGemma29BFree = "google/gemma-2-9b-it:free"

// GoogleGemma29B is the ID for model Google: Gemma 2 9B
//
// Gemma 2 9B by Google is an advanced, open-source language model that sets a new standard for efficiency and performance in its size class.
//
// Designed for a wide variety of tasks, it empowers developers and researchers to build innovative applications, while maintaining accessibility, safety, and cost-effectiveness.
//
// See the [launch announcement](https://blog.google/technology/developers/google-gemma-2/) for more details. Usage of Gemma is subject to Google's [Gemma Terms of Use](https://ai.google.dev/gemma/terms).
const GoogleGemma29B = "google/gemma-2-9b-it"

// 01AIYiLarge is the ID for model 01.AI: Yi Large
//
// The Yi Large model was designed by 01.AI with the following usecases in mind: knowledge search, data classification, human-like chat bots, and customer service.
//
// It stands out for its multilingual proficiency, particularly in Spanish, Chinese, Japanese, German, and French.
//
// Check out the [launch announcement](https://01-ai.github.io/blog/01.ai-yi-large-llm-launch) to learn more.
const O1AIYiLarge = "01-ai/yi-large"

// AI21JambaInstruct is the ID for model AI21: Jamba Instruct
//
// The Jamba-Instruct model, introduced by AI21 Labs, is an instruction-tuned variant of their hybrid SSM-Transformer Jamba model, specifically optimized for enterprise applications.
//
// - 256K Context Window: It can process extensive information, equivalent to a 400-page novel, which is beneficial for tasks involving large documents such as financial reports or legal documents
// - Safety and Accuracy: Jamba-Instruct is designed with enhanced safety features to ensure secure deployment in enterprise environments, reducing the risk and cost of implementation
//
// Read their [announcement](https://www.ai21.com/blog/announcing-jamba) to learn more.
//
// Jamba has a knowledge cutoff of February 2024.
const AI21JambaInstruct = "ai21/jamba-instruct"

// AnthropicClaude35Sonnet20240620SelfModerated is the ID for model Anthropic: Claude 3.5 Sonnet (2024-06-20) (self-moderated)
//
// Claude 3.5 Sonnet delivers better-than-Opus capabilities, faster-than-Sonnet speeds, at the same Sonnet prices. Sonnet is particularly good at:
//
// - Coding: Autonomously writes, edits, and runs code with reasoning and troubleshooting
// - Data science: Augments human data science expertise; navigates unstructured data while using multiple tools for insights
// - Visual processing: excelling at interpreting charts, graphs, and images, accurately transcribing text to derive insights beyond just the text alone
// - Agentic tasks: exceptional tool use, making it great at agentic tasks (i.e. complex, multi-step problem solving tasks that require engaging with other systems)
//
// For the latest version (2024-10-23), check out [Claude 3.5 Sonnet](/anthropic/claude-3.5-sonnet).
//
// #multimodal
const AnthropicClaude35Sonnet20240620SelfModerated = "anthropic/claude-3.5-sonnet-20240620:beta"

// AnthropicClaude35Sonnet20240620 is the ID for model Anthropic: Claude 3.5 Sonnet (2024-06-20)
//
// Claude 3.5 Sonnet delivers better-than-Opus capabilities, faster-than-Sonnet speeds, at the same Sonnet prices. Sonnet is particularly good at:
//
// - Coding: Autonomously writes, edits, and runs code with reasoning and troubleshooting
// - Data science: Augments human data science expertise; navigates unstructured data while using multiple tools for insights
// - Visual processing: excelling at interpreting charts, graphs, and images, accurately transcribing text to derive insights beyond just the text alone
// - Agentic tasks: exceptional tool use, making it great at agentic tasks (i.e. complex, multi-step problem solving tasks that require engaging with other systems)
//
// For the latest version (2024-10-23), check out [Claude 3.5 Sonnet](/anthropic/claude-3.5-sonnet).
//
// #multimodal
const AnthropicClaude35Sonnet20240620 = "anthropic/claude-3.5-sonnet-20240620"

// Sao10kLlama3Euryale70BV21 is the ID for model Sao10k: Llama 3 Euryale 70B v2.1
//
// Euryale 70B v2.1 is a model focused on creative roleplay from [Sao10k](https://ko-fi.com/sao10k).
//
// - Better prompt adherence.
// - Better anatomy / spatial awareness.
// - Adapts much better to unique and custom formatting / reply formats.
// - Very creative, lots of unique swipes.
// - Is not restrictive during roleplays.
const Sao10kLlama3Euryale70BV21 = "sao10k/l3-euryale-70b"

// Dolphin292Mixtral8x22B is the ID for model Dolphin 2.9.2 Mixtral 8x22B 🐬
//
// Dolphin 2.9 is designed for instruction following, conversational, and coding. This model is a finetune of [Mixtral 8x22B Instruct](/models/mistralai/mixtral-8x22b-instruct). It features a 64k context length and was fine-tuned with a 16k sequence length using ChatML templates.
//
// This model is a successor to [Dolphin Mixtral 8x7B](/models/cognitivecomputations/dolphin-mixtral-8x7b).
//
// The model is uncensored and is stripped of alignment and bias. It requires an external alignment layer for ethical use. Users are cautioned to use this highly compliant model responsibly, as detailed in a blog post about uncensored models at [erichartford.com/uncensored-models](https://erichartford.com/uncensored-models).
//
// #moe #uncensored
const Dolphin292Mixtral8x22B = "cognitivecomputations/dolphin-mixtral-8x22b"

// Qwen272BInstruct is the ID for model Qwen 2 72B Instruct
//
// Qwen2 72B is a transformer-based model that excels in language understanding, multilingual capabilities, coding, mathematics, and reasoning.
//
// It features SwiGLU activation, attention QKV bias, and group query attention. It is pretrained on extensive data with supervised finetuning and direct preference optimization.
//
// For more details, see this [blog post](https://qwenlm.github.io/blog/qwen2/) and [GitHub repo](https://github.com/QwenLM/Qwen2).
//
// Usage of this model is subject to [Tongyi Qianwen LICENSE AGREEMENT](https://huggingface.co/Qwen/Qwen1.5-110B-Chat/blob/main/LICENSE).
const Qwen272BInstruct = "qwen/qwen-2-72b-instruct"

// MistralMistral7BInstructFree is the ID for model Mistral: Mistral 7B Instruct (free)
//
// A high-performing, industry-standard 7.3B parameter model, with optimizations for speed and context length.
//
// *Mistral 7B Instruct has multiple version variants, and this is intended to be the latest version.*
const MistralMistral7BInstructFree = "mistralai/mistral-7b-instruct:free"

// MistralMistral7BInstruct is the ID for model Mistral: Mistral 7B Instruct
//
// A high-performing, industry-standard 7.3B parameter model, with optimizations for speed and context length.
//
// *Mistral 7B Instruct has multiple version variants, and this is intended to be the latest version.*
const MistralMistral7BInstruct = "mistralai/mistral-7b-instruct"

// MistralMistral7BInstructV03 is the ID for model Mistral: Mistral 7B Instruct v0.3
//
// A high-performing, industry-standard 7.3B parameter model, with optimizations for speed and context length.
//
// An improved version of [Mistral 7B Instruct v0.2](/models/mistralai/mistral-7b-instruct-v0.2), with the following changes:
//
// - Extended vocabulary to 32768
// - Supports v3 Tokenizer
// - Supports function calling
//
// NOTE: Support for function calling depends on the provider.
const MistralMistral7BInstructV03 = "mistralai/mistral-7b-instruct-v0.3"

// NousResearchHermes2ProLlama38B is the ID for model NousResearch: Hermes 2 Pro - Llama-3 8B
//
// Hermes 2 Pro is an upgraded, retrained version of Nous Hermes 2, consisting of an updated and cleaned version of the OpenHermes 2.5 Dataset, as well as a newly introduced Function Calling and JSON Mode dataset developed in-house.
const NousResearchHermes2ProLlama38B = "nousresearch/hermes-2-pro-llama-3-8b"

// MicrosoftPhi3Mini128KInstructFree is the ID for model Microsoft: Phi-3 Mini 128K Instruct (free)
//
// Phi-3 Mini is a powerful 3.8B parameter model designed for advanced language understanding, reasoning, and instruction following. Optimized through supervised fine-tuning and preference adjustments, it excels in tasks involving common sense, mathematics, logical reasoning, and code processing.
//
// At time of release, Phi-3 Medium demonstrated state-of-the-art performance among lightweight models. This model is static, trained on an offline dataset with an October 2023 cutoff date.
const MicrosoftPhi3Mini128KInstructFree = "microsoft/phi-3-mini-128k-instruct:free"

// MicrosoftPhi3Mini128KInstruct is the ID for model Microsoft: Phi-3 Mini 128K Instruct
//
// Phi-3 Mini is a powerful 3.8B parameter model designed for advanced language understanding, reasoning, and instruction following. Optimized through supervised fine-tuning and preference adjustments, it excels in tasks involving common sense, mathematics, logical reasoning, and code processing.
//
// At time of release, Phi-3 Medium demonstrated state-of-the-art performance among lightweight models. This model is static, trained on an offline dataset with an October 2023 cutoff date.
const MicrosoftPhi3Mini128KInstruct = "microsoft/phi-3-mini-128k-instruct"

// MicrosoftPhi3Medium128KInstructFree is the ID for model Microsoft: Phi-3 Medium 128K Instruct (free)
//
// Phi-3 128K Medium is a powerful 14-billion parameter model designed for advanced language understanding, reasoning, and instruction following. Optimized through supervised fine-tuning and preference adjustments, it excels in tasks involving common sense, mathematics, logical reasoning, and code processing.
//
// At time of release, Phi-3 Medium demonstrated state-of-the-art performance among lightweight models. In the MMLU-Pro eval, the model even comes close to a Llama3 70B level of performance.
//
// For 4k context length, try [Phi-3 Medium 4K](/models/microsoft/phi-3-medium-4k-instruct).
const MicrosoftPhi3Medium128KInstructFree = "microsoft/phi-3-medium-128k-instruct:free"

// MicrosoftPhi3Medium128KInstruct is the ID for model Microsoft: Phi-3 Medium 128K Instruct
//
// Phi-3 128K Medium is a powerful 14-billion parameter model designed for advanced language understanding, reasoning, and instruction following. Optimized through supervised fine-tuning and preference adjustments, it excels in tasks involving common sense, mathematics, logical reasoning, and code processing.
//
// At time of release, Phi-3 Medium demonstrated state-of-the-art performance among lightweight models. In the MMLU-Pro eval, the model even comes close to a Llama3 70B level of performance.
//
// For 4k context length, try [Phi-3 Medium 4K](/models/microsoft/phi-3-medium-4k-instruct).
const MicrosoftPhi3Medium128KInstruct = "microsoft/phi-3-medium-128k-instruct"

// NeverSleepLlama3Lumimaid70B is the ID for model NeverSleep: Llama 3 Lumimaid 70B
//
// The NeverSleep team is back, with a Llama 3 70B finetune trained on their curated roleplay data. Striking a balance between eRP and RP, Lumimaid was designed to be serious, yet uncensored when necessary.
//
// To enhance it's overall intelligence and chat capability, roughly 40% of the training data was not roleplay. This provides a breadth of knowledge to access, while still keeping roleplay as the primary strength.
//
// Usage of this model is subject to [Meta's Acceptable Use Policy](https://llama.meta.com/llama3/use-policy/).
const NeverSleepLlama3Lumimaid70B = "neversleep/llama-3-lumimaid-70b"

// GoogleGeminiFlash15 is the ID for model Google: Gemini Flash 1.5
//
// Gemini 1.5 Flash is a foundation model that performs well at a variety of multimodal tasks such as visual understanding, classification, summarization, and creating content from image, audio and video. It's adept at processing visual and text inputs such as photographs, documents, infographics, and screenshots.
//
// Gemini 1.5 Flash is designed for high-volume, high-frequency tasks where cost and latency matter. On most common tasks, Flash achieves comparable quality to other Gemini Pro models at a significantly reduced cost. Flash is well-suited for applications like chat assistants and on-demand content generation where speed and scale matter.
//
// Usage of Gemini is subject to Google's [Gemini Terms of Use](https://ai.google.dev/terms).
//
// #multimodal
const GoogleGeminiFlash15 = "google/gemini-flash-1.5"

// DeepSeekV25 is the ID for model DeepSeek V2.5
//
// DeepSeek-V2.5 is an upgraded version that combines DeepSeek-V2-Chat and DeepSeek-Coder-V2-Instruct. The new model integrates the general and coding abilities of the two previous versions. For model details, please visit [DeepSeek-V2 page](https://github.com/deepseek-ai/DeepSeek-V2) for more information.
const DeepSeekV25 = "deepseek/deepseek-chat-v2.5"

// OpenAIGPT4o20240513 is the ID for model OpenAI: GPT-4o (2024-05-13)
//
// GPT-4o ("o" for "omni") is OpenAI's latest AI model, supporting both text and image inputs with text outputs. It maintains the intelligence level of [GPT-4 Turbo](/models/openai/gpt-4-turbo) while being twice as fast and 50% more cost-effective. GPT-4o also offers improved performance in processing non-English languages and enhanced visual capabilities.
//
// For benchmarking against other models, it was briefly called ["im-also-a-good-gpt2-chatbot"](https://twitter.com/LiamFedus/status/1790064963966370209)
//
// #multimodal
const OpenAIGPT4o20240513 = "openai/gpt-4o-2024-05-13"

// MetaLlamaGuard28B is the ID for model Meta: LlamaGuard 2 8B
//
// This safeguard model has 8B parameters and is based on the Llama 3 family. Just like is predecessor, [LlamaGuard 1](https://huggingface.co/meta-llama/LlamaGuard-7b), it can do both prompt and response classification.
//
// LlamaGuard 2 acts as a normal LLM would, generating text that indicates whether the given input/output is safe/unsafe. If deemed unsafe, it will also share the content categories violated.
//
// For best results, please use raw prompt input or the `/completions` endpoint, instead of the chat API.
//
// It has demonstrated strong performance compared to leading closed-source models in human evaluations.
//
// To read more about the model release, [click here](https://ai.meta.com/blog/meta-llama-3/). Usage of this model is subject to [Meta's Acceptable Use Policy](https://llama.meta.com/llama3/use-policy/).
const MetaLlamaGuard28B = "meta-llama/llama-guard-2-8b"

// OpenAIGPT4o is the ID for model OpenAI: GPT-4o
//
// GPT-4o ("o" for "omni") is OpenAI's latest AI model, supporting both text and image inputs with text outputs. It maintains the intelligence level of [GPT-4 Turbo](/models/openai/gpt-4-turbo) while being twice as fast and 50% more cost-effective. GPT-4o also offers improved performance in processing non-English languages and enhanced visual capabilities.
//
// For benchmarking against other models, it was briefly called ["im-also-a-good-gpt2-chatbot"](https://twitter.com/LiamFedus/status/1790064963966370209)
//
// #multimodal
const OpenAIGPT4o = "openai/gpt-4o"

// OpenAIGPT4oExtended is the ID for model OpenAI: GPT-4o (extended)
//
// GPT-4o ("o" for "omni") is OpenAI's latest AI model, supporting both text and image inputs with text outputs. It maintains the intelligence level of [GPT-4 Turbo](/models/openai/gpt-4-turbo) while being twice as fast and 50% more cost-effective. GPT-4o also offers improved performance in processing non-English languages and enhanced visual capabilities.
//
// For benchmarking against other models, it was briefly called ["im-also-a-good-gpt2-chatbot"](https://twitter.com/LiamFedus/status/1790064963966370209)
//
// #multimodal
const OpenAIGPT4oExtended = "openai/gpt-4o:extended"

// NeverSleepLlama3Lumimaid8BExtended is the ID for model NeverSleep: Llama 3 Lumimaid 8B (extended)
//
// The NeverSleep team is back, with a Llama 3 8B finetune trained on their curated roleplay data. Striking a balance between eRP and RP, Lumimaid was designed to be serious, yet uncensored when necessary.
//
// To enhance it's overall intelligence and chat capability, roughly 40% of the training data was not roleplay. This provides a breadth of knowledge to access, while still keeping roleplay as the primary strength.
//
// Usage of this model is subject to [Meta's Acceptable Use Policy](https://llama.meta.com/llama3/use-policy/).
const NeverSleepLlama3Lumimaid8BExtended = "neversleep/llama-3-lumimaid-8b:extended"

// NeverSleepLlama3Lumimaid8B is the ID for model NeverSleep: Llama 3 Lumimaid 8B
//
// The NeverSleep team is back, with a Llama 3 8B finetune trained on their curated roleplay data. Striking a balance between eRP and RP, Lumimaid was designed to be serious, yet uncensored when necessary.
//
// To enhance it's overall intelligence and chat capability, roughly 40% of the training data was not roleplay. This provides a breadth of knowledge to access, while still keeping roleplay as the primary strength.
//
// Usage of this model is subject to [Meta's Acceptable Use Policy](https://llama.meta.com/llama3/use-policy/).
const NeverSleepLlama3Lumimaid8B = "neversleep/llama-3-lumimaid-8b"

// Fimbulvetr11BV2 is the ID for model Fimbulvetr 11B v2
//
// Creative writing model, routed with permission. It's fast, it keeps the conversation going, and it stays in character.
//
// If you submit a raw prompt, you can use Alpaca or Vicuna formats.
const Fimbulvetr11BV2 = "sao10k/fimbulvetr-11b-v2"

// MetaLlama38BInstructFree is the ID for model Meta: Llama 3 8B Instruct (free)
//
// Meta's latest class of model (Llama 3) launched with a variety of sizes & flavors. This 8B instruct-tuned version was optimized for high quality dialogue usecases.
//
// It has demonstrated strong performance compared to leading closed-source models in human evaluations.
//
// To read more about the model release, [click here](https://ai.meta.com/blog/meta-llama-3/). Usage of this model is subject to [Meta's Acceptable Use Policy](https://llama.meta.com/llama3/use-policy/).
const MetaLlama38BInstructFree = "meta-llama/llama-3-8b-instruct:free"

// MetaLlama38BInstruct is the ID for model Meta: Llama 3 8B Instruct
//
// Meta's latest class of model (Llama 3) launched with a variety of sizes & flavors. This 8B instruct-tuned version was optimized for high quality dialogue usecases.
//
// It has demonstrated strong performance compared to leading closed-source models in human evaluations.
//
// To read more about the model release, [click here](https://ai.meta.com/blog/meta-llama-3/). Usage of this model is subject to [Meta's Acceptable Use Policy](https://llama.meta.com/llama3/use-policy/).
const MetaLlama38BInstruct = "meta-llama/llama-3-8b-instruct"

// MetaLlama370BInstruct is the ID for model Meta: Llama 3 70B Instruct
//
// Meta's latest class of model (Llama 3) launched with a variety of sizes & flavors. This 70B instruct-tuned version was optimized for high quality dialogue usecases.
//
// It has demonstrated strong performance compared to leading closed-source models in human evaluations.
//
// To read more about the model release, [click here](https://ai.meta.com/blog/meta-llama-3/). Usage of this model is subject to [Meta's Acceptable Use Policy](https://llama.meta.com/llama3/use-policy/).
const MetaLlama370BInstruct = "meta-llama/llama-3-70b-instruct"

// MistralMixtral8x22BInstruct is the ID for model Mistral: Mixtral 8x22B Instruct
//
// Mistral's official instruct fine-tuned version of [Mixtral 8x22B](/models/mistralai/mixtral-8x22b). It uses 39B active parameters out of 141B, offering unparalleled cost efficiency for its size. Its strengths include:
// - strong math, coding, and reasoning
// - large context length (64k)
// - fluency in English, French, Italian, German, and Spanish
//
// See benchmarks on the launch announcement [here](https://mistral.ai/news/mixtral-8x22b/).
// #moe
const MistralMixtral8x22BInstruct = "mistralai/mixtral-8x22b-instruct"

// WizardLM28x22B is the ID for model WizardLM-2 8x22B
//
// WizardLM-2 8x22B is Microsoft AI's most advanced Wizard model. It demonstrates highly competitive performance compared to leading proprietary models, and it consistently outperforms all existing state-of-the-art opensource models.
//
// It is an instruct finetune of [Mixtral 8x22B](/models/mistralai/mixtral-8x22b).
//
// To read more about the model release, [click here](https://wizardlm.github.io/WizardLM2/).
//
// #moe
const WizardLM28x22B = "microsoft/wizardlm-2-8x22b"

// WizardLM27B is the ID for model WizardLM-2 7B
//
// WizardLM-2 7B is the smaller variant of Microsoft AI's latest Wizard model. It is the fastest and achieves comparable performance with existing 10x larger opensource leading models
//
// It is a finetune of [Mistral 7B Instruct](/models/mistralai/mistral-7b-instruct), using the same technique as [WizardLM-2 8x22B](/models/microsoft/wizardlm-2-8x22b).
//
// To read more about the model release, [click here](https://wizardlm.github.io/WizardLM2/).
//
// #moe
const WizardLM27B = "microsoft/wizardlm-2-7b"

// GoogleGeminiPro15 is the ID for model Google: Gemini Pro 1.5
//
// Google's latest multimodal model, supports image and video[0] in text or chat prompts.
//
// Optimized for language tasks including:
//
// - Code generation
// - Text generation
// - Text editing
// - Problem solving
// - Recommendations
// - Information extraction
// - Data extraction or generation
// - AI agents
//
// Usage of Gemini is subject to Google's [Gemini Terms of Use](https://ai.google.dev/terms).
//
// * [0]: Video input is not available through OpenRouter at this time.
const GoogleGeminiPro15 = "google/gemini-pro-1.5"

// OpenAIGPT4Turbo is the ID for model OpenAI: GPT-4 Turbo
//
// The latest GPT-4 Turbo model with vision capabilities. Vision requests can now use JSON mode and function calling.
//
// Training data: up to December 2023.
const OpenAIGPT4Turbo = "openai/gpt-4-turbo"

// CohereCommandR is the ID for model Cohere: Command R+
//
// Command R+ is a new, 104B-parameter LLM from Cohere. It's useful for roleplay, general consumer usecases, and Retrieval Augmented Generation (RAG).
//
// It offers multilingual support for ten key languages to facilitate global business operations. See benchmarks and the launch post [here](https://txt.cohere.com/command-r-plus-microsoft-azure/).
//
// Use of this model is subject to Cohere's [Usage Policy](https://docs.cohere.com/docs/usage-policy) and [SaaS Agreement](https://cohere.com/saas-agreement).
const CohereCommandRPlus = "cohere/command-r-plus"

// CohereCommandR042024 is the ID for model Cohere: Command R+ (04-2024)
//
// Command R+ is a new, 104B-parameter LLM from Cohere. It's useful for roleplay, general consumer usecases, and Retrieval Augmented Generation (RAG).
//
// It offers multilingual support for ten key languages to facilitate global business operations. See benchmarks and the launch post [here](https://txt.cohere.com/command-r-plus-microsoft-azure/).
//
// Use of this model is subject to Cohere's [Usage Policy](https://docs.cohere.com/docs/usage-policy) and [SaaS Agreement](https://cohere.com/saas-agreement).
const CohereCommandR042024 = "cohere/command-r-plus-04-2024"

// DatabricksDBRX132BInstruct is the ID for model Databricks: DBRX 132B Instruct
//
// DBRX is a new open source large language model developed by Databricks. At 132B, it outperforms existing open source LLMs like Llama 2 70B and [Mixtral-8x7b](/models/mistralai/mixtral-8x7b) on standard industry benchmarks for language understanding, programming, math, and logic.
//
// It uses a fine-grained mixture-of-experts (MoE) architecture. 36B parameters are active on any input. It was pre-trained on 12T tokens of text and code data. Compared to other open MoE models like Mixtral-8x7B and Grok-1, DBRX is fine-grained, meaning it uses a larger number of smaller experts.
//
// See the launch announcement and benchmark results [here](https://www.databricks.com/blog/introducing-dbrx-new-state-art-open-llm).
//
// #moe
const DatabricksDBRX132BInstruct = "databricks/dbrx-instruct"

// MidnightRose70B is the ID for model Midnight Rose 70B
//
// A merge with a complex family tree, this model was crafted for roleplaying and storytelling. Midnight Rose is a successor to Rogue Rose and Aurora Nights and improves upon them both. It wants to produce lengthy output by default and is the best creative writing merge produced so far by sophosympatheia.
//
// Descending from earlier versions of Midnight Rose and [Wizard Tulu Dolphin 70B](https://huggingface.co/sophosympatheia/Wizard-Tulu-Dolphin-70B-v1.0), it inherits the best qualities of each.
const MidnightRose70B = "sophosympatheia/midnight-rose-70b"

// CohereCommand is the ID for model Cohere: Command
//
// Command is an instruction-following conversational model that performs language tasks with high quality, more reliably and with a longer context than our base generative models.
//
// Use of this model is subject to Cohere's [Usage Policy](https://docs.cohere.com/docs/usage-policy) and [SaaS Agreement](https://cohere.com/saas-agreement).
const CohereCommand = "cohere/command"

// CohereCommandR is the ID for model Cohere: Command R
//
// Command-R is a 35B parameter model that performs conversational language tasks at a higher quality, more reliably, and with a longer context than previous models. It can be used for complex workflows like code generation, retrieval augmented generation (RAG), tool use, and agents.
//
// Read the launch post [here](https://txt.cohere.com/command-r/).
//
// Use of this model is subject to Cohere's [Usage Policy](https://docs.cohere.com/docs/usage-policy) and [SaaS Agreement](https://cohere.com/saas-agreement).
const CohereCommandR = "cohere/command-r"

// AnthropicClaude3HaikuSelfModerated is the ID for model Anthropic: Claude 3 Haiku (self-moderated)
//
// Claude 3 Haiku is Anthropic's fastest and most compact model for
// near-instant responsiveness. Quick and accurate targeted performance.
//
// See the launch announcement and benchmark results [here](https://www.anthropic.com/news/claude-3-haiku)
//
// #multimodal
const AnthropicClaude3HaikuSelfModerated = "anthropic/claude-3-haiku:beta"

// AnthropicClaude3Haiku is the ID for model Anthropic: Claude 3 Haiku
//
// Claude 3 Haiku is Anthropic's fastest and most compact model for
// near-instant responsiveness. Quick and accurate targeted performance.
//
// See the launch announcement and benchmark results [here](https://www.anthropic.com/news/claude-3-haiku)
//
// #multimodal
const AnthropicClaude3Haiku = "anthropic/claude-3-haiku"

// AnthropicClaude3OpusSelfModerated is the ID for model Anthropic: Claude 3 Opus (self-moderated)
//
// Claude 3 Opus is Anthropic's most powerful model for highly complex tasks. It boasts top-level performance, intelligence, fluency, and understanding.
//
// See the launch announcement and benchmark results [here](https://www.anthropic.com/news/claude-3-family)
//
// #multimodal
const AnthropicClaude3OpusSelfModerated = "anthropic/claude-3-opus:beta"

// AnthropicClaude3Opus is the ID for model Anthropic: Claude 3 Opus
//
// Claude 3 Opus is Anthropic's most powerful model for highly complex tasks. It boasts top-level performance, intelligence, fluency, and understanding.
//
// See the launch announcement and benchmark results [here](https://www.anthropic.com/news/claude-3-family)
//
// #multimodal
const AnthropicClaude3Opus = "anthropic/claude-3-opus"

// AnthropicClaude3SonnetSelfModerated is the ID for model Anthropic: Claude 3 Sonnet (self-moderated)
//
// Claude 3 Sonnet is an ideal balance of intelligence and speed for enterprise workloads. Maximum utility at a lower price, dependable, balanced for scaled deployments.
//
// See the launch announcement and benchmark results [here](https://www.anthropic.com/news/claude-3-family)
//
// #multimodal
const AnthropicClaude3SonnetSelfModerated = "anthropic/claude-3-sonnet:beta"

// AnthropicClaude3Sonnet is the ID for model Anthropic: Claude 3 Sonnet
//
// Claude 3 Sonnet is an ideal balance of intelligence and speed for enterprise workloads. Maximum utility at a lower price, dependable, balanced for scaled deployments.
//
// See the launch announcement and benchmark results [here](https://www.anthropic.com/news/claude-3-family)
//
// #multimodal
const AnthropicClaude3Sonnet = "anthropic/claude-3-sonnet"

// CohereCommandR032024 is the ID for model Cohere: Command R (03-2024)
//
// Command-R is a 35B parameter model that performs conversational language tasks at a higher quality, more reliably, and with a longer context than previous models. It can be used for complex workflows like code generation, retrieval augmented generation (RAG), tool use, and agents.
//
// Read the launch post [here](https://txt.cohere.com/command-r/).
//
// Use of this model is subject to Cohere's [Usage Policy](https://docs.cohere.com/docs/usage-policy) and [SaaS Agreement](https://cohere.com/saas-agreement).
const CohereCommandR032024 = "cohere/command-r-03-2024"

// MistralLarge is the ID for model Mistral Large
//
// This is Mistral AI's flagship model, Mistral Large 2 (version `mistral-large-2407`). It's a proprietary weights-available model and excels at reasoning, code, JSON, chat, and more. Read the launch announcement [here](https://mistral.ai/news/mistral-large-2407/).
//
// It supports dozens of languages including French, German, Spanish, Italian, Portuguese, Arabic, Hindi, Russian, Chinese, Japanese, and Korean, along with 80+ coding languages including Python, Java, C, C++, JavaScript, and Bash. Its long context window allows precise information recall from large documents.
const MistralLarge = "mistralai/mistral-large"

// GoogleGemma7B is the ID for model Google: Gemma 7B
//
// Gemma by Google is an advanced, open-source language model family, leveraging the latest in decoder-only, text-to-text technology. It offers English language capabilities across text generation tasks like question answering, summarization, and reasoning. The Gemma 7B variant is comparable in performance to leading open source models.
//
// Usage of Gemma is subject to Google's [Gemma Terms of Use](https://ai.google.dev/gemma/terms).
const GoogleGemma7B = "google/gemma-7b-it"

// OpenAIGPT35TurboOlderV0613 is the ID for model OpenAI: GPT-3.5 Turbo (older v0613)
//
// GPT-3.5 Turbo is OpenAI's fastest model. It can understand and generate natural language or code, and is optimized for chat and traditional completion tasks.
//
// Training data up to Sep 2021.
const OpenAIGPT35TurboOlderV0613 = "openai/gpt-3.5-turbo-0613"

// OpenAIGPT4TurboPreview is the ID for model OpenAI: GPT-4 Turbo Preview
//
// The preview GPT-4 model with improved instruction following, JSON mode, reproducible outputs, parallel function calling, and more. Training data: up to Dec 2023.
//
// **Note:** heavily rate limited by OpenAI while in preview.
const OpenAIGPT4TurboPreview = "openai/gpt-4-turbo-preview"

// NousHermes2Mixtral8x7BDPO is the ID for model Nous: Hermes 2 Mixtral 8x7B DPO
//
// Nous Hermes 2 Mixtral 8x7B DPO is the new flagship Nous Research model trained over the [Mixtral 8x7B MoE LLM](/models/mistralai/mixtral-8x7b).
//
// The model was trained on over 1,000,000 entries of primarily [GPT-4](/models/openai/gpt-4) generated data, as well as other high quality data from open datasets across the AI landscape, achieving state of the art performance on a variety of tasks.
//
// #moe
const NousHermes2Mixtral8x7BDPO = "nousresearch/nous-hermes-2-mixtral-8x7b-dpo"

// MistralSmall is the ID for model Mistral Small
//
// With 22 billion parameters, Mistral Small v24.09 offers a convenient mid-point between (Mistral NeMo 12B)[/mistralai/mistral-nemo] and (Mistral Large 2)[/mistralai/mistral-large], providing a cost-effective solution that can be deployed across various platforms and environments. It has better reasoning, exhibits more capabilities, can produce and reason about code, and is multiligual, supporting English, French, German, Italian, and Spanish.
const MistralSmall = "mistralai/mistral-small"

// MistralTiny is the ID for model Mistral Tiny
//
// This model is currently powered by Mistral-7B-v0.2, and incorporates a "better" fine-tuning than [Mistral 7B](/models/mistralai/mistral-7b-instruct-v0.1), inspired by community work. It's best used for large batch processing tasks where cost is a significant factor but reasoning capabilities are not crucial.
const MistralTiny = "mistralai/mistral-tiny"

// MistralMedium is the ID for model Mistral Medium
//
// This is Mistral AI's closed-source, medium-sided model. It's powered by a closed-source prototype and excels at reasoning, code, JSON, chat, and more. In benchmarks, it compares with many of the flagship models of other companies.
const MistralMedium = "mistralai/mistral-medium"

// Dolphin26Mixtral8x7B is the ID for model Dolphin 2.6 Mixtral 8x7B 🐬
//
// This is a 16k context fine-tune of [Mixtral-8x7b](/models/mistralai/mixtral-8x7b). It excels in coding tasks due to extensive training with coding data and is known for its obedience, although it lacks DPO tuning.
//
// The model is uncensored and is stripped of alignment and bias. It requires an external alignment layer for ethical use. Users are cautioned to use this highly compliant model responsibly, as detailed in a blog post about uncensored models at [erichartford.com/uncensored-models](https://erichartford.com/uncensored-models).
//
// #moe #uncensored
const Dolphin26Mixtral8x7B = "cognitivecomputations/dolphin-mixtral-8x7b"

// GoogleGeminiProVision10 is the ID for model Google: Gemini Pro Vision 1.0
//
// Google's flagship multimodal model, supporting image and video in text or chat prompts for a text or code response.
//
// See the benchmarks and prompting guidelines from [Deepmind](https://deepmind.google/technologies/gemini/).
//
// Usage of Gemini is subject to Google's [Gemini Terms of Use](https://ai.google.dev/terms).
//
// #multimodal
const GoogleGeminiProVision10 = "google/gemini-pro-vision"

// GoogleGeminiPro10 is the ID for model Google: Gemini Pro 1.0
//
// Google's flagship text generation model. Designed to handle natural language tasks, multiturn text and code chat, and code generation.
//
// See the benchmarks and prompting guidelines from [Deepmind](https://deepmind.google/technologies/gemini/).
//
// Usage of Gemini is subject to Google's [Gemini Terms of Use](https://ai.google.dev/terms).
const GoogleGeminiPro10 = "google/gemini-pro"

// MistralMixtral8x7BBase is the ID for model Mistral: Mixtral 8x7B (base)
//
// Mixtral 8x7B is a pretrained generative Sparse Mixture of Experts, by Mistral AI. Incorporates 8 experts (feed-forward networks) for a total of 47B parameters. Base model (not fine-tuned for instructions) - see [Mixtral 8x7B Instruct](/models/mistralai/mixtral-8x7b-instruct) for an instruct-tuned model.
//
// #moe
const MistralMixtral8x7BBase = "mistralai/mixtral-8x7b"

// MistralMixtral8x7BInstruct is the ID for model Mistral: Mixtral 8x7B Instruct
//
// Mixtral 8x7B Instruct is a pretrained generative Sparse Mixture of Experts, by Mistral AI, for chat and instruction use. Incorporates 8 experts (feed-forward networks) for a total of 47 billion parameters.
//
// Instruct model fine-tuned by Mistral. #moe
const MistralMixtral8x7BInstruct = "mistralai/mixtral-8x7b-instruct"

// OpenChat357BFree is the ID for model OpenChat 3.5 7B (free)
//
// OpenChat 7B is a library of open-source language models, fine-tuned with "C-RLFT (Conditioned Reinforcement Learning Fine-Tuning)" - a strategy inspired by offline reinforcement learning. It has been trained on mixed-quality data without preference labels.
//
// - For OpenChat fine-tuned on Mistral 7B, check out [OpenChat 7B](/models/openchat/openchat-7b).
// - For OpenChat fine-tuned on Llama 8B, check out [OpenChat 8B](/models/openchat/openchat-8b).
//
// #open-source
const OpenChat357BFree = "openchat/openchat-7b:free"

// OpenChat357B is the ID for model OpenChat 3.5 7B
//
// OpenChat 7B is a library of open-source language models, fine-tuned with "C-RLFT (Conditioned Reinforcement Learning Fine-Tuning)" - a strategy inspired by offline reinforcement learning. It has been trained on mixed-quality data without preference labels.
//
// - For OpenChat fine-tuned on Mistral 7B, check out [OpenChat 7B](/models/openchat/openchat-7b).
// - For OpenChat fine-tuned on Llama 8B, check out [OpenChat 8B](/models/openchat/openchat-8b).
//
// #open-source
const OpenChat357B = "openchat/openchat-7b"

// Noromaid20B is the ID for model Noromaid 20B
//
// A collab between IkariDev and Undi. This merge is suitable for RP, ERP, and general knowledge.
//
// #merge #uncensored
const Noromaid20B = "neversleep/noromaid-20b"

// AnthropicClaudeV2SelfModerated is the ID for model Anthropic: Claude v2 (self-moderated)
//
// Claude 2 delivers advancements in key capabilities for enterprises—including an industry-leading 200K token context window, significant reductions in rates of model hallucination, system prompts and a new beta feature: tool use.
const AnthropicClaudeV2SelfModerated = "anthropic/claude-2:beta"

// AnthropicClaudeV2 is the ID for model Anthropic: Claude v2
//
// Claude 2 delivers advancements in key capabilities for enterprises—including an industry-leading 200K token context window, significant reductions in rates of model hallucination, system prompts and a new beta feature: tool use.
const AnthropicClaudeV2 = "anthropic/claude-2"

// AnthropicClaudeV21SelfModerated is the ID for model Anthropic: Claude v2.1 (self-moderated)
//
// Claude 2 delivers advancements in key capabilities for enterprises—including an industry-leading 200K token context window, significant reductions in rates of model hallucination, system prompts and a new beta feature: tool use.
const AnthropicClaudeV21SelfModerated = "anthropic/claude-2.1:beta"

// AnthropicClaudeV21 is the ID for model Anthropic: Claude v2.1
//
// Claude 2 delivers advancements in key capabilities for enterprises—including an industry-leading 200K token context window, significant reductions in rates of model hallucination, system prompts and a new beta feature: tool use.
const AnthropicClaudeV21 = "anthropic/claude-2.1"

// OpenHermes25Mistral7B is the ID for model OpenHermes 2.5 Mistral 7B
//
// A continuation of [OpenHermes 2 model](/models/teknium/openhermes-2-mistral-7b), trained on additional code datasets.
// Potentially the most interesting finding from training on a good ratio (est. of around 7-14% of the total dataset) of code instruction was that it has boosted several non-code benchmarks, including TruthfulQA, AGIEval, and GPT4All suite. It did however reduce BigBench benchmark score, but the net gain overall is significant.
const OpenHermes25Mistral7B = "teknium/openhermes-2.5-mistral-7b"

// ToppyM7BFree is the ID for model Toppy M 7B (free)
//
// A wild 7B parameter model that merges several models using the new task_arithmetic merge method from mergekit.
// List of merged models:
// - NousResearch/Nous-Capybara-7B-V1.9
// - [HuggingFaceH4/zephyr-7b-beta](/models/huggingfaceh4/zephyr-7b-beta)
// - lemonilia/AshhLimaRP-Mistral-7B
// - Vulkane/120-Days-of-Sodom-LoRA-Mistral-7b
// - Undi95/Mistral-pippa-sharegpt-7b-qlora
//
// #merge #uncensored
const ToppyM7BFree = "undi95/toppy-m-7b:free"

// ToppyM7B is the ID for model Toppy M 7B
//
// A wild 7B parameter model that merges several models using the new task_arithmetic merge method from mergekit.
// List of merged models:
// - NousResearch/Nous-Capybara-7B-V1.9
// - [HuggingFaceH4/zephyr-7b-beta](/models/huggingfaceh4/zephyr-7b-beta)
// - lemonilia/AshhLimaRP-Mistral-7B
// - Vulkane/120-Days-of-Sodom-LoRA-Mistral-7b
// - Undi95/Mistral-pippa-sharegpt-7b-qlora
//
// #merge #uncensored
const ToppyM7B = "undi95/toppy-m-7b"

// Goliath120B is the ID for model Goliath 120B
//
// A large LLM created by combining two fine-tuned Llama 70B models into one 120B model. Combines Xwin and Euryale.
//
// Credits to
// - [@chargoddard](https://huggingface.co/chargoddard) for developing the framework used to merge the model - [mergekit](https://github.com/cg123/mergekit).
// - [@Undi95](https://huggingface.co/Undi95) for helping with the merge ratios.
//
// #merge
const Goliath120B = "alpindale/goliath-120b"

// AutoRouter is the ID for model Auto Router
//
// Your prompt will be processed by a meta-model and routed to one of dozens of models (see below), optimizing for the best possible output.
//
// To see which model was used, visit [Activity](/activity), or read the `model` attribute of the response. Your response will be priced at the same rate as the routed model.
//
// The meta-model is powered by [Not Diamond](https://docs.notdiamond.ai/docs/how-not-diamond-works). Learn more in our [docs](/docs/model-routing).
//
// Requests will be routed to the following models:
// - [openai/gpt-4o-2024-08-06](/openai/gpt-4o-2024-08-06)
// - [openai/gpt-4o-2024-05-13](/openai/gpt-4o-2024-05-13)
// - [openai/gpt-4o-mini-2024-07-18](/openai/gpt-4o-mini-2024-07-18)
// - [openai/chatgpt-4o-latest](/openai/chatgpt-4o-latest)
// - [openai/o1-preview-2024-09-12](/openai/o1-preview-2024-09-12)
// - [openai/o1-mini-2024-09-12](/openai/o1-mini-2024-09-12)
// - [anthropic/claude-3.5-sonnet](/anthropic/claude-3.5-sonnet)
// - [anthropic/claude-3.5-haiku](/anthropic/claude-3.5-haiku)
// - [anthropic/claude-3-opus](/anthropic/claude-3-opus)
// - [anthropic/claude-2.1](/anthropic/claude-2.1)
// - [google/gemini-pro-1.5](/google/gemini-pro-1.5)
// - [google/gemini-flash-1.5](/google/gemini-flash-1.5)
// - [mistralai/mistral-large-2407](/mistralai/mistral-large-2407)
// - [mistralai/mistral-nemo](/mistralai/mistral-nemo)
// - [deepseek/deepseek-r1](/deepseek/deepseek-r1)
// - [meta-llama/llama-3.1-70b-instruct](/meta-llama/llama-3.1-70b-instruct)
// - [meta-llama/llama-3.1-405b-instruct](/meta-llama/llama-3.1-405b-instruct)
// - [mistralai/mixtral-8x22b-instruct](/mistralai/mixtral-8x22b-instruct)
// - [cohere/command-r-plus](/cohere/command-r-plus)
// - [cohere/command-r](/cohere/command-r)
const AutoRouter = "openrouter/auto"

// OpenAIGPT35Turbo16kOlderV1106 is the ID for model OpenAI: GPT-3.5 Turbo 16k (older v1106)
//
// An older GPT-3.5 Turbo model with improved instruction following, JSON mode, reproducible outputs, parallel function calling, and more. Training data: up to Sep 2021.
const OpenAIGPT35Turbo16kOlderV1106 = "openai/gpt-3.5-turbo-1106"

// OpenAIGPT4TurboOlderV1106 is the ID for model OpenAI: GPT-4 Turbo (older v1106)
//
// The latest GPT-4 Turbo model with vision capabilities. Vision requests can now use JSON mode and function calling.
//
// Training data: up to April 2023.
const OpenAIGPT4TurboOlderV1106 = "openai/gpt-4-1106-preview"

// GooglePaLM2Chat32k is the ID for model Google: PaLM 2 Chat 32k
//
// PaLM 2 is a language model by Google with improved multilingual, reasoning and coding capabilities.
const GooglePaLM2Chat32k = "google/palm-2-chat-bison-32k"

// GooglePaLM2CodeChat32k is the ID for model Google: PaLM 2 Code Chat 32k
//
// PaLM 2 fine-tuned for chatbot conversations that help with code-related questions.
const GooglePaLM2CodeChat32k = "google/palm-2-codechat-bison-32k"

// Airoboros70B is the ID for model Airoboros 70B
//
// A Llama 2 70B fine-tune using synthetic data (the Airoboros dataset).
//
// Currently based on [jondurbin/airoboros-l2-70b](https://huggingface.co/jondurbin/airoboros-l2-70b-2.2.1), but might get updated in the future.
const Airoboros70B = "jondurbin/airoboros-l2-70b"

// Xwin70B is the ID for model Xwin 70B
//
// Xwin-LM aims to develop and open-source alignment tech for LLMs. Our first release, built-upon on the [Llama2](/models/${Model.Llama_2_13B_Chat}) base models, ranked TOP-1 on AlpacaEval. Notably, it's the first to surpass [GPT-4](/models/${Model.GPT_4}) on this benchmark. The project will be continuously updated.
const Xwin70B = "xwin-lm/xwin-lm-70b"

// OpenAIGPT35TurboInstruct is the ID for model OpenAI: GPT-3.5 Turbo Instruct
//
// This model is a variant of GPT-3.5 Turbo tuned for instructional prompts and omitting chat-related optimizations. Training data: up to Sep 2021.
const OpenAIGPT35TurboInstruct = "openai/gpt-3.5-turbo-instruct"

// MistralMistral7BInstructV01 is the ID for model Mistral: Mistral 7B Instruct v0.1
//
// A 7.3B parameter model that outperforms Llama 2 13B on all benchmarks, with optimizations for speed and context length.
const MistralMistral7BInstructV01 = "mistralai/mistral-7b-instruct-v0.1"

// PygmalionMythalion13B is the ID for model Pygmalion: Mythalion 13B
//
// A blend of the new Pygmalion-13b and MythoMax. #merge
const PygmalionMythalion13B = "pygmalionai/mythalion-13b"

// OpenAIGPT35Turbo16k is the ID for model OpenAI: GPT-3.5 Turbo 16k
//
// This model offers four times the context length of gpt-3.5-turbo, allowing it to support approximately 20 pages of text in a single request at a higher cost. Training data: up to Sep 2021.
const OpenAIGPT35Turbo16k = "openai/gpt-3.5-turbo-16k"

// OpenAIGPT432k is the ID for model OpenAI: GPT-4 32k
//
// GPT-4-32k is an extended version of GPT-4, with the same capabilities but quadrupled context length, allowing for processing up to 40 pages of text in a single pass. This is particularly beneficial for handling longer content like interacting with PDFs without an external vector database. Training data: up to Sep 2021.
const OpenAIGPT432k = "openai/gpt-4-32k"

// OpenAIGPT432kOlderV0314 is the ID for model OpenAI: GPT-4 32k (older v0314)
//
// GPT-4-32k is an extended version of GPT-4, with the same capabilities but quadrupled context length, allowing for processing up to 40 pages of text in a single pass. This is particularly beneficial for handling longer content like interacting with PDFs without an external vector database. Training data: up to Sep 2021.
const OpenAIGPT432kOlderV0314 = "openai/gpt-4-32k-0314"

// NousHermes13B is the ID for model Nous: Hermes 13B
//
// A state-of-the-art language model fine-tuned on over 300k instructions by Nous Research, with Teknium and Emozilla leading the fine tuning process.
const NousHermes13B = "nousresearch/nous-hermes-llama2-13b"

// MancerWeaverAlpha is the ID for model Mancer: Weaver (alpha)
//
// An attempt to recreate Claude-style verbosity, but don't expect the same level of coherence or memory. Meant for use in roleplay/narrative situations.
const MancerWeaverAlpha = "mancer/weaver"

// HuggingFaceZephyr7BFree is the ID for model Hugging Face: Zephyr 7B (free)
//
// Zephyr is a series of language models that are trained to act as helpful assistants. Zephyr-7B-β is the second model in the series, and is a fine-tuned version of [mistralai/Mistral-7B-v0.1](/models/mistralai/mistral-7b-instruct-v0.1) that was trained on a mix of publicly available, synthetic datasets using Direct Preference Optimization (DPO).
const HuggingFaceZephyr7BFree = "huggingfaceh4/zephyr-7b-beta:free"

// AnthropicClaudeV20SelfModerated is the ID for model Anthropic: Claude v2.0 (self-moderated)
//
// Anthropic's flagship model. Superior performance on tasks that require complex reasoning. Supports hundreds of pages of text.
const AnthropicClaudeV20SelfModerated = "anthropic/claude-2.0:beta"

// AnthropicClaudeV20 is the ID for model Anthropic: Claude v2.0
//
// Anthropic's flagship model. Superior performance on tasks that require complex reasoning. Supports hundreds of pages of text.
const AnthropicClaudeV20 = "anthropic/claude-2.0"

// ReMMSLERP13B is the ID for model ReMM SLERP 13B
//
// A recreation trial of the original MythoMax-L2-B13 but with updated models. #merge
const ReMMSLERP13B = "undi95/remm-slerp-l2-13b"

// GooglePaLM2Chat is the ID for model Google: PaLM 2 Chat
//
// PaLM 2 is a language model by Google with improved multilingual, reasoning and coding capabilities.
const GooglePaLM2Chat = "google/palm-2-chat-bison"

// GooglePaLM2CodeChat is the ID for model Google: PaLM 2 Code Chat
//
// PaLM 2 fine-tuned for chatbot conversations that help with code-related questions.
const GooglePaLM2CodeChat = "google/palm-2-codechat-bison"

// MythoMax13BFree is the ID for model MythoMax 13B (free)
//
// One of the highest performing and most popular fine-tunes of Llama 2 13B, with rich descriptions and roleplay. #merge
const MythoMax13BFree = "gryphe/mythomax-l2-13b:free"

// MythoMax13B is the ID for model MythoMax 13B
//
// One of the highest performing and most popular fine-tunes of Llama 2 13B, with rich descriptions and roleplay. #merge
const MythoMax13B = "gryphe/mythomax-l2-13b"

// MetaLlama213BChat is the ID for model Meta: Llama 2 13B Chat
//
// A 13 billion parameter language model from Meta, fine tuned for chat completions
const MetaLlama213BChat = "meta-llama/llama-2-13b-chat"

// MetaLlama270BChat is the ID for model Meta: Llama 2 70B Chat
//
// The flagship, 70 billion parameter language model from Meta, fine tuned for chat completions. Llama 2 is an auto-regressive language model that uses an optimized transformer architecture. The tuned versions use supervised fine-tuning (SFT) and reinforcement learning with human feedback (RLHF) to align to human preferences for helpfulness and safety.
const MetaLlama270BChat = "meta-llama/llama-2-70b-chat"

// OpenAIGPT35Turbo is the ID for model OpenAI: GPT-3.5 Turbo
//
// GPT-3.5 Turbo is OpenAI's fastest model. It can understand and generate natural language or code, and is optimized for chat and traditional completion tasks.
//
// Training data up to Sep 2021.
const OpenAIGPT35Turbo = "openai/gpt-3.5-turbo"

// OpenAIGPT35Turbo16k is the ID for model OpenAI: GPT-3.5 Turbo 16k
//
// The latest GPT-3.5 Turbo model with improved instruction following, JSON mode, reproducible outputs, parallel function calling, and more. Training data: up to Sep 2021.
//
// This version has a higher accuracy at responding in requested formats and a fix for a bug which caused a text encoding issue for non-English language function calls.
const OpenAIGPT35Turbo0125 = "openai/gpt-3.5-turbo-0125"

// OpenAIGPT4 is the ID for model OpenAI: GPT-4
//
// OpenAI's flagship model, GPT-4 is a large-scale multimodal language model capable of solving difficult problems with greater accuracy than previous models due to its broader general knowledge and advanced reasoning capabilities. Training data: up to Sep 2021.
const OpenAIGPT4 = "openai/gpt-4"

// OpenAIGPT4OlderV0314 is the ID for model OpenAI: GPT-4 (older v0314)
//
// GPT-4-0314 is the first version of GPT-4 released, with a context length of 8,192 tokens, and was supported until June 14. Training data: up to Sep 2021.
const OpenAIGPT4OlderV0314 = "openai/gpt-4-0314"
