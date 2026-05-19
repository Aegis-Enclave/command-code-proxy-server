package proxy

import (
	"encoding/json"
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
					Type:  "tool_use",
					ID:    strPtr(tc.ID),
					Name:  strPtr(tc.Function.Name),
					Input: parseToolInput(tc.Function.Arguments),
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

func ConvertTools(openAITools []any) []any {
	if len(openAITools) == 0 {
		return []any{}
	}

	tools := make([]any, 0, len(openAITools))
	for _, tool := range openAITools {
		toolMap, ok := tool.(map[string]any)
		if !ok {
			continue
		}

		toolType, _ := toolMap["type"].(string)
		if toolType != "function" {
			tools = append(tools, toolMap)
			continue
		}

		fn, ok := toolMap["function"].(map[string]any)
		if !ok {
			continue
		}

		name, _ := fn["name"].(string)
		if name == "" {
			continue
		}

		inputSchema, ok := fn["parameters"].(map[string]any)
		if !ok || inputSchema == nil {
			inputSchema = map[string]any{"type": "object", "properties": map[string]any{}}
		}

		ccTool := map[string]any{
			"name":         name,
			"input_schema": inputSchema,
		}
		if description, ok := fn["description"].(string); ok && description != "" {
			ccTool["description"] = description
		}
		tools = append(tools, ccTool)
	}

	return tools
}

func parseToolInput(arguments string) any {
	if arguments == "" {
		return map[string]any{}
	}
	var input any
	if err := json.Unmarshal([]byte(arguments), &input); err != nil {
		return map[string]any{"arguments": arguments}
	}
	return input
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

func contentPartToString(content any) string {
	switch v := content.(type) {
	case string:
		return v
	case []any:
		var b strings.Builder
		for _, item := range v {
			if m, ok := item.(map[string]any); ok {
				if text, ok := m["text"].(string); ok {
					b.WriteString(text)
				}
			}
		}
		return b.String()
	default:
		data, err := json.Marshal(v)
		if err != nil {
			return ""
		}
		return string(data)
	}
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
			partMap, ok := part.(map[string]any)
			if !ok {
				continue
			}
			typ, _ := partMap["type"].(string)
			switch typ {
			case "text":
				if text, ok := partMap["text"].(string); ok && text != "" {
					parts = append(parts, api.CCContentPart{Type: "text", Text: strPtr(text)})
				}
			case "image_url":
				if imgURL, ok := partMap["image_url"].(map[string]any); ok {
					if url, ok := imgURL["url"].(string); ok && url != "" {
						parts = append(parts, api.CCContentPart{Type: "text", Text: strPtr("[Image URL: " + url + "]")})
					}
				}
			case "tool_use":
				id, _ := partMap["id"].(string)
				name, _ := partMap["name"].(string)
				input := partMap["input"]
				parts = append(parts, api.CCContentPart{Type: "tool_use", ID: strPtr(id), Name: strPtr(name), Input: input})
			case "tool-call":
				id, _ := partMap["id"].(string)
				name, _ := partMap["name"].(string)
				input := partMap["input"]
				if input == nil {
					input = partMap["arguments"]
				}
				parts = append(parts, api.CCContentPart{Type: "tool-call", ID: strPtr(id), Name: strPtr(name), Input: input})
			case "tool_result", "tool-result":
				toolID, _ := partMap["tool_use_id"].(string)
				if toolID == "" {
					toolID, _ = partMap["toolCallId"].(string)
				}
				toolName, _ := partMap["toolName"].(string)
				parts = append(parts, api.CCContentPart{Type: "tool-result", ToolCallID: strPtr(toolID), ToolName: strPtr(toolName), Text: strPtr(contentPartToString(partMap["content"]))})
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
