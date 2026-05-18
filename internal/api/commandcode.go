package api

// CommandCode API types (internal)

type CCContentPart struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type CCMessage struct {
	Role    string          `json:"role"`
	Content []CCContentPart `json:"content"`
}

type CCChatParams struct {
	Model       string      `json:"model"`
	Messages    []CCMessage `json:"messages"`
	Tools       []any       `json:"tools"`
	System      string      `json:"system"`
	MaxTokens   int         `json:"max_tokens"`
	Temperature float64     `json:"temperature"`
	Stream      bool        `json:"stream"`
}

type CCConfig struct {
	WorkingDir    string   `json:"workingDir"`
	Date          string   `json:"date"`
	Environment   string   `json:"environment"`
	Structure     []string `json:"structure"`
	IsGitRepo     bool     `json:"isGitRepo"`
	CurrentBranch string   `json:"currentBranch"`
	MainBranch    string   `json:"mainBranch"`
	GitStatus     string   `json:"gitStatus"`
	RecentCommits []string `json:"recentCommits"`
}

type CCRequestBody struct {
	Config   CCConfig     `json:"config"`
	Memory   string       `json:"memory"`
	Taste    string       `json:"taste"`
	Skills   string       `json:"skills"`
	Params   CCChatParams `json:"params"`
	ThreadID string       `json:"threadId"`
}

type CCStreamEvent struct {
	Type         string `json:"type"`
	Text         string `json:"text"`
	FinishReason string `json:"finishReason"`
	Error        *struct {
		Message    string `json:"message"`
		StatusCode *int   `json:"statusCode"`
	} `json:"error"`
	TotalUsage *struct {
		InputTokens  int `json:"inputTokens"`
		OutputTokens int `json:"outputTokens"`
	} `json:"totalUsage"`
}
