package codec

import (
	"fmt"

	"go.jetify.com/ai/api"
)

var (
	ProviderName = "openai"

	FileSearchToolName = "file_search"
	FileSearchToolID   = fmt.Sprintf("%s.%s", ProviderName, FileSearchToolName)

	ComputerUseToolName = "computer_use_preview"
	ComputerUseToolID   = fmt.Sprintf("%s.%s", ProviderName, ComputerUseToolName)

	WebSearchToolName = "web_search_preview"
	WebSearchToolID   = fmt.Sprintf("%s.%s", ProviderName, WebSearchToolName)
)

// FileSearchToolArgs is a built-in tool that searches for relevant content from uploaded files.
// Learn more about the [file search tool](https://platform.openai.com/docs/guides/tools-file-search).
type FileSearchToolArgs struct {
	// The IDs of the vector stores to search.
	VectorStoreIDs []string `json:"vector_store_ids,omitzero"`

	// The maximum number of results to return. This number should be between 1 and 50
	// inclusive. If not provided, it's set to a default.
	MaxNumResults int `json:"max_num_results,omitzero"`

	// TODO: Add filters and ranking options
	// // A filter to apply based on file attributes.
	// Filters X `json:"filters,omitzero"`
	// // Ranking options for search.
	// RankingOptions X `json:"ranking_options,omitzero"`
}

// FileSearchToolCall represents the results of a file search operation.
// See the [file search guide](https://platform.openai.com/docs/guides/tools-file-search)
// for more information.
type FileSearchToolCall struct {
	// Queries contains the search terms used to find files
	Queries []string `json:"queries"`
	// Results holds the matching files after executing the file search.
	Results []FileSearchResult `json:"results"`
}

// FileSearchResult contains metadata and content for a single file match.
type FileSearchResult struct {
	// FileID uniquely identifies the file
	FileID string `json:"file_id"`
	// Filename is the name of the matched file
	Filename string `json:"filename"`
	// Score indicates the relevance of the match (0.0 to 1.0)
	Score float64 `json:"score"`
	// Text contains the retrieved file content
	Text string `json:"text"`
}

// FileSearchToolOption allows customizing file search tool configuration.
type FileSearchToolOption func(*FileSearchToolArgs)

// WithVectorStoreIDs sets the vector store IDs to search.
func WithVectorStoreIDs(ids ...string) FileSearchToolOption {
	return func(args *FileSearchToolArgs) {
		args.VectorStoreIDs = ids
	}
}

// WithMaxNumResults sets the maximum number of results to return.
func WithMaxNumResults(maxResults int) FileSearchToolOption {
	return func(args *FileSearchToolArgs) {
		args.MaxNumResults = maxResults
	}
}

// FileSearchTool creates a new file search tool with the specified configuration.
// FileSearchTool is a built-in tool that searches for relevant content from uploaded files.
// Learn more about the [file search tool](https://platform.openai.com/docs/guides/tools-file-search).
func FileSearchTool(options ...FileSearchToolOption) *api.ProviderDefinedTool {
	args := &FileSearchToolArgs{}

	// Apply options
	for _, opt := range options {
		opt(args)
	}

	return &api.ProviderDefinedTool{
		ID:   FileSearchToolID,
		Name: FileSearchToolName,
		Args: args,
	}
}

// WebSearchToolArgs is a built-in tool that searches the web for relevant results to use in a response.
// Learn more about the [web search tool](https://platform.openai.com/docs/guides/tools-web-search).
type WebSearchToolArgs struct {
	// High level guidance for the amount of context window space to use for the
	// search. One of `low`, `medium`, or `high`. `medium` is the default.
	SearchContextSize string `json:"search_context_size,omitempty"`
	// User location information for geographically relevant results
	UserLocation *WebSearchUserLocation `json:"user_location,omitempty"`
}

// WebSearchUserLocation represents the user location information for a web search
type WebSearchUserLocation struct {
	// Free text input for the city of the user, e.g. `San Francisco`.
	City string `json:"city,omitzero"`
	// The two-letter [ISO country code](https://en.wikipedia.org/wiki/ISO_3166-1) of
	// the user, e.g. `US`.
	Country string `json:"country,omitzero"`
	// Free text input for the region of the user, e.g. `California`.
	Region string `json:"region,omitzero"`
	// The [IANA timezone](https://timeapi.io/documentation/iana-timezones) of the
	// user, e.g. `America/Los_Angeles`.
	Timezone string `json:"timezone,omitzero"`
}

// WebSearchToolOption allows customizing web search tool configuration.
type WebSearchToolOption func(*WebSearchToolArgs)

// WithSearchContextSize sets the search context size.
func WithSearchContextSize(size string) WebSearchToolOption {
	return func(args *WebSearchToolArgs) {
		args.SearchContextSize = size
	}
}

// WithUserLocation sets the user location for geographically relevant results.
func WithUserLocation(location *WebSearchUserLocation) WebSearchToolOption {
	return func(args *WebSearchToolArgs) {
		args.UserLocation = location
	}
}

// WebSearchTool creates a new web search tool with the specified configuration.
// WebSearchTool is a built-in tool that searches the web for relevant results to use in a response.
// Learn more about the [web search tool](https://platform.openai.com/docs/guides/tools-web-search).
func WebSearchTool(options ...WebSearchToolOption) *api.ProviderDefinedTool {
	args := &WebSearchToolArgs{}

	// Apply options
	for _, opt := range options {
		opt(args)
	}

	return &api.ProviderDefinedTool{
		ID:   WebSearchToolID,
		Name: WebSearchToolName,
		Args: args,
	}
}

// ComputerUseToolArgs is a built-in tool that controls a virtual computer. Learn more about the
// [computer tool](https://platform.openai.com/docs/guides/tools-computer-use).
//
// The properties DisplayHeight, DisplayWidth, Environment, Type are required.
type ComputerUseToolArgs struct {
	// The height of the computer display.
	DisplayHeight int `json:"display_height,omitempty"`
	// The width of the computer display.
	DisplayWidth int `json:"display_width,omitempty"`
	// The type of computer environment to control.
	//
	// Any of "mac", "windows", "ubuntu", "browser".
	Environment string `json:"environment,omitempty"`
}

// ComputerToolCall represents a computer-based tool operation.  See the
// [computer use guide](https://platform.openai.com/docs/guides/tools-computer-use)
// for more information.
type ComputerToolCall struct {
	// Action represents the type of action to perform. Any of "click", "double_click",
	// "drag", "keypress", "move", "screenshot", "scroll", "type", or "wait".
	Action string

	// Coordinates represents the screen coordinates to perform the action on, if
	// applicable.
	// Applies to these types of actions: "click", "double_click", "move".
	Coordinates ComputerCoordinates

	// MouseButton indicates which mouse button to press for a "click" action. One of
	// "left", "right", "wheel", "back", or "forward". Assume "left" if not specified.
	MouseButton string

	// DragPath is the path of coordinates to follow for a "drag" action. Coordinates
	// will appear as an array of coordinate objects, eg
	//
	// ```
	// [
	//
	//	{ x: 100, y: 200 },
	//	{ x: 200, y: 300 }
	//
	// ]
	// ```
	DragPath []ComputerCoordinates

	// Keys indicates the combination of keys the model is requesting to be pressed
	// for a "keypress" action. This is an array of strings, each representing a
	// key to be pressed simultaneously.
	Keys []string

	// ScrollDistance indicates the distance to scroll in the x and y directions for
	// a "scroll" action.
	ScrollDistance ComputerCoordinates

	// Text is the text that should be typed in a "type" action.
	Text string
}

type ComputerCoordinates struct {
	X int
	Y int
}

// ComputerSafetyCheck represents a pending safety check for the computer call.
//
// The properties ID, Code, Message are required.
type ComputerSafetyCheck struct {
	// The ID of the pending safety check.
	ID string
	// The type of the pending safety check.
	Code string
	// Details about the pending safety check.
	Message string
}

// ComputerUseToolOption allows customizing computer use tool configuration.
type ComputerUseToolOption func(*ComputerUseToolArgs)

// WithDisplaySize sets the display dimensions.
func WithDisplaySize(width, height int) ComputerUseToolOption {
	return func(args *ComputerUseToolArgs) {
		args.DisplayWidth = width
		args.DisplayHeight = height
	}
}

// WithEnvironment sets the computer environment to control.
func WithEnvironment(env string) ComputerUseToolOption {
	return func(args *ComputerUseToolArgs) {
		args.Environment = env
	}
}

// ComputerUseTool creates a new computer use tool with the specified configuration.
// ComputerUseTool is a built-in tool that controls a virtual computer. Learn more about the
// [computer tool](https://platform.openai.com/docs/guides/tools-computer-use).
func ComputerUseTool(displayWidth, displayHeight int, environment string, options ...ComputerUseToolOption) *api.ProviderDefinedTool {
	args := &ComputerUseToolArgs{
		DisplayWidth:  displayWidth,
		DisplayHeight: displayHeight,
		Environment:   environment,
	}

	// Apply options
	for _, opt := range options {
		opt(args)
	}

	return &api.ProviderDefinedTool{
		ID:   ComputerUseToolID,
		Name: ComputerUseToolName,
		Args: args,
	}
}
