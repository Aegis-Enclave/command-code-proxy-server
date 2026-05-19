package api

// OpenAI-compatible types (client-facing)

type OpenAIMessage struct {
	Role          string        `json:"role"`
	Content       interface{}   `json:"content,omitempty"`
	Name          string        `json:"name,omitempty"`
	ToolCalls     []ToolCall    `json:"tool_calls,omitempty"`
	ToolCallID    string        `json:"tool_call_id,omitempty"`
	Refusal       string        `json:"refusal,omitempty"`
	Audio         *MessageAudio `json:"audio,omitempty"`
}

type ContentPart struct {
	Type     string    `json:"type"`
	Text     string    `json:"text,omitempty"`
	ImageURL *ImageURL `json:"image_url,omitempty"`
}

type ImageURL struct {
	URL       string `json:"url"`
	Detail    string `json:"detail,omitempty"`
	Modalities string `json:"modalities,omitempty"`
}

type ToolCall struct {
	ID       string      `json:"id"`
	Type     string      `json:"type"`
	Function FunctionCall `json:"function"`
}

type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type MessageAudio struct {
	ID       string `json:"id"`
	Data     string `json:"data"`
	Duration float64 `json:"duration"`
}

type OpenAIChatRequest struct {
	Model       string          `json:"model"`
	Messages    []OpenAIMessage `json:"messages"`
	Temperature *float64        `json:"temperature,omitempty"`
	MaxTokens   *int            `json:"max_tokens,omitempty"`
	Stream      bool            `json:"stream,omitempty"`
	Tools       []any           `json:"tools,omitempty"`
}

type OpenAIChoice struct {
	Index        int           `json:"index"`
	Message      *OpenAIMessage `json:"message,omitempty"`
	Delta        *OpenAIDelta  `json:"delta,omitempty"`
	FinishReason *string       `json:"finish_reason,omitempty"`
}

type OpenAIDelta struct {
	Role      string            `json:"role,omitempty"`
	Content   string            `json:"content,omitempty"`
	ToolCalls []OpenAIDeltaToolCall `json:"tool_calls,omitempty"`
}

type OpenAIDeltaToolCall struct {
	Index    int    `json:"index"`
	ID       string `json:"id,omitempty"`
	Type     string `json:"type,omitempty"`
	Function *OpenAIDeltaFunction `json:"function,omitempty"`
}

type OpenAIDeltaFunction struct {
	Name      string `json:"name,omitempty"`
	Arguments string `json:"arguments,omitempty"`
}

type OpenAIUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type OpenAIChatResponse struct {
	ID      string         `json:"id"`
	Object  string         `json:"object"`
	Created int64          `json:"created"`
	Model   string         `json:"model"`
	Choices []OpenAIChoice `json:"choices"`
	Usage   *OpenAIUsage   `json:"usage,omitempty"`
}

type OpenAIModel struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	OwnedBy string `json:"owned_by"`
}

type OpenAIModelList struct {
	Object string        `json:"object"`
	Data   []OpenAIModel `json:"data"`
}
