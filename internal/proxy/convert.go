package proxy

import (
	"strings"

	"github.com/dev2k6/command-code-proxy-server/internal/api"
)

// Convert OpenAI messages to CommandCode format
func ConvertMessages(openAIMsgs []api.OpenAIMessage) []api.CCMessage {
	var ccMsgs []api.CCMessage
	for _, m := range openAIMsgs {
		// Convert tool role to user with tool-result type
		if m.Role == "tool" {
			ccMsgs = append(ccMsgs, api.CCMessage{
				Role: "user",
				Content: []api.CCContentPart{{
					Type:       "tool-result",
					ToolCallID: strPtr(m.ToolCallID),
					ToolName:   strPtr(m.Name),
					Text:       strPtr(contentToString(m.Content)),
				}},
			})
			continue
		}

		// Convert assistant with tool_calls
		if m.Role == "assistant" && len(m.ToolCalls) > 0 {
			contentParts := parseContent(m.Content)
			for _, tc := range m.ToolCalls {
				contentParts = append(contentParts, api.CCContentPart{
					Type:       "tool_use",
					ToolCallID: strPtr(tc.ID),
					ToolName:   strPtr(tc.Function.Name),
					Text:       strPtr(tc.Function.Arguments),
				})
			}
			ccMsgs = append(ccMsgs, api.CCMessage{
				Role:    m.Role,
				Content: contentParts,
			})
			continue
		}

		contentParts := parseContent(m.Content)
		ccMsgs = append(ccMsgs, api.CCMessage{
			Role:    m.Role,
			Content: contentParts,
		})
	}
	return ccMsgs
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func contentToString(content interface{}) string {
	switch v := content.(type) {
	case string:
		return v
	case []any:
		for _, part := range v {
			if partMap, ok := part.(map[string]any); ok {
				if text, ok := partMap["text"].(string); ok {
					return text
				}
			}
		}
	}
	return ""
}

func parseContent(content interface{}) []api.CCContentPart {
	switch v := content.(type) {
	case string:
		if v == "" {
			return nil
		}
		return []api.CCContentPart{{Type: "text", Text: strPtr(v)}}
	case []any:
		var parts []api.CCContentPart
		for _, part := range v {
			if partMap, ok := part.(map[string]any); ok {
				if typ, ok := partMap["type"].(string); ok {
					p := api.CCContentPart{Type: typ}
					if text, ok := partMap["text"].(string); ok {
						p.Text = strPtr(text)
					}
					if imgURL, ok := partMap["image_url"].(map[string]any); ok {
						if url, ok := imgURL["url"].(string); ok {
							merged := ""
							if p.Text != nil {
								merged = *p.Text
							}
							merged = merged + "\n[Image URL: " + url + "]"
							p.Text = strPtr(merged)
						}
					}
					parts = append(parts, p)
				}
			}
		}
		return parts
	default:
		return nil
	}
}

// Extract system message and remaining messages
func ExtractSystem(msgs []api.OpenAIMessage) (string, []api.OpenAIMessage) {
	var system strings.Builder
	var rest []api.OpenAIMessage
	for _, m := range msgs {
		if m.Role == "system" {
			if system.Len() > 0 {
				system.WriteString("\n")
			}
			if contentStr, ok := m.Content.(string); ok {
				system.WriteString(contentStr)
			}
		} else {
			rest = append(rest, m)
		}
	}
	return system.String(), rest
}
