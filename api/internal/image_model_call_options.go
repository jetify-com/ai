package api

import "net/http"

// ImageCallOptions represents the options for generating images.
type ImageCallOptions struct {
	// N is the number of images to generate.
	N int

	// Size of the images to generate.
	// Must have the format `{width}x{height}`.
	// nil will use the provider's default size.
	Size *string

	// AspectRatio of the images to generate.
	// Must have the format `{width}:{height}`.
	// nil will use the provider's default aspect ratio.
	AspectRatio *string

	// Seed for the image generation.
	// nil will use the provider's default seed.
	Seed *int

	// ProviderOptions are additional provider-specific options that are passed through to the provider
	// as body parameters.
	//
	// The outer map is keyed by the provider name, and the inner
	// map is keyed by the provider-specific metadata key. The value can be any JSON-compatible value
	// (string, number, boolean, null, array, or object).
	// Example:
	//   {
	//     "openai": {
	//       "style": "vivid",
	//       "quality": 1,
	//       "hd": true,
	//       "metadata": {
	//         "user": "test"
	//       }
	//     }
	//   }
	ProviderOptions map[string]map[string]any

	// Headers are additional HTTP headers to be sent with the request.
	// Only applicable for HTTP-based providers.
	Headers http.Header
}

// ImageCallOption is a function that modifies ImageCallOptions.
type ImageCallOption func(*ImageCallOptions)

// WithImageCount sets the number of images to generate.
// N is the number of images to generate.
func WithImageCount(n int) ImageCallOption {
	return func(o *ImageCallOptions) {
		o.N = n
	}
}

// WithImageSize sets the size of images to generate.
// Must have the format `{width}x{height}`.
// nil will use the provider's default size.
func WithImageSize(size string) ImageCallOption {
	return func(o *ImageCallOptions) {
		o.Size = &size
	}
}

// WithImageAspectRatio sets the aspect ratio of images to generate.
// Must have the format `{width}:{height}`.
// nil will use the provider's default aspect ratio.
func WithImageAspectRatio(ratio string) ImageCallOption {
	return func(o *ImageCallOptions) {
		o.AspectRatio = &ratio
	}
}

// WithImageSeed sets the seed for image generation.
// nil will use the provider's default seed.
func WithImageSeed(seed int) ImageCallOption {
	return func(o *ImageCallOptions) {
		o.Seed = &seed
	}
}

// WithImageProviderOptions sets provider-specific options that are passed through to the provider
// as body parameters.
//
// The outer map is keyed by the provider name, and the inner
// map is keyed by the provider-specific metadata key. The value can be any JSON-compatible value
// (string, number, boolean, null, array, or object).
// Example:
//
//	{
//	  "openai": {
//	    "style": "vivid",
//	    "quality": 1,
//	    "hd": true,
//	    "metadata": {
//	      "user": "test"
//	    }
//	  }
//	}
func WithImageProviderOptions(options map[string]map[string]any) ImageCallOption {
	return func(o *ImageCallOptions) {
		o.ProviderOptions = options
	}
}

// WithImageHeaders sets additional HTTP headers to be sent with the request.
// Only applicable for HTTP-based providers.
func WithImageHeaders(headers http.Header) ImageCallOption {
	return func(o *ImageCallOptions) {
		o.Headers = headers
	}
}
