package proxy

import (
	"strings"

	"github.com/dev2k6/command-code-proxy-server/internal/api"
)

// Convert OpenAI messages to CommandCode format
func ConvertMessages(openAIMsgs []api.OpenAIMessage) []api.CCMessage {
	var ccMsgs []api.CCMessage
	for _, m := range openAIMsgs {
		contentParts := parseContent(m.Content)
		ccMsgs = append(ccMsgs, api.CCMessage{
			Role:    m.Role,
			Content: contentParts,
		})
	}
	return ccMsgs
}

func parseContent(content interface{}) []api.CCContentPart {
	switch v := content.(type) {
	case string:
		if v == "" {
			return nil
		}
		return []api.CCContentPart{{Type: "text", Text: v}}
	case []any:
		var parts []api.CCContentPart
		for _, part := range v {
			if partMap, ok := part.(map[string]any); ok {
				if typ, ok := partMap["type"].(string); ok {
					p := api.CCContentPart{Type: typ}
					if text, ok := partMap["text"].(string); ok {
						p.Text = text
					}
					// Note: CommandCode API may not support image_url
					// For now, we skip image parts
					if imgURL, ok := partMap["image_url"].(map[string]any); ok {
						if url, ok := imgURL["url"].(string); ok {
							// Try to extract text from image (e.g., base64 data)
							// or just include the URL as text
							p.Text = p.Text + "\n[Image URL: " + url + "]"
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
