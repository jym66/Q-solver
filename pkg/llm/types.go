package llm

// ChunkType 流式输出块类型
type ChunkType string

const (
	ChunkThinking ChunkType = "thinking" // 思维链
	ChunkContent  ChunkType = "content"  // 正文内容
)

// StreamChunk 统一的流式输出块
type StreamChunk struct {
	Type    ChunkType `json:"type"`
	Content string    `json:"content"`
}

// Role 消息角色
type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

// ContentType 内容类型
type ContentType string

const (
	ContentText  ContentType = "text"
	ContentImage ContentType = "image"
	ContentPDF   ContentType = "pdf"
)

// ContentPart 内容块（支持文本、图片、PDF）
type ContentPart struct {
	Type   ContentType `json:"type"`
	Text   string      `json:"text,omitempty"`
	Base64 string      `json:"base64,omitempty"` // 包含 data:xxx;base64, 前缀
}

// Message 统一的消息格式
type Message struct {
	Role     Role          `json:"role"`
	Content  string        `json:"content,omitempty"`  // 纯文本内容
	Parts    []ContentPart `json:"parts,omitempty"`    // 多模态内容
	Thinking string        `json:"thinking,omitempty"` // 思维链（仅 Assistant）
}

// StreamCallback 统一的流式回调
type StreamCallback func(chunk StreamChunk)

// NewTextMessage 创建纯文本消息
func NewTextMessage(role Role, content string) Message {
	return Message{Role: role, Content: content}
}

// NewSystemMessage 创建系统消息
func NewSystemMessage(content string) Message {
	return NewTextMessage(RoleSystem, content)
}

// NewUserMessage 创建用户消息
func NewUserMessage(content string) Message {
	return NewTextMessage(RoleUser, content)
}

// NewAssistantMessage 创建助手消息
func NewAssistantMessage(content string) Message {
	return NewTextMessage(RoleAssistant, content)
}

// NewMultiPartMessage 创建多模态消息
func NewMultiPartMessage(role Role, parts []ContentPart) Message {
	return Message{Role: role, Parts: parts}
}

// TextPart 创建文本内容块
func TextPart(text string) ContentPart {
	return ContentPart{Type: ContentText, Text: text}
}

// ImagePart 创建图片内容块
func ImagePart(base64URL string) ContentPart {
	return ContentPart{Type: ContentImage, Base64: base64URL}
}

// PDFPart 创建 PDF 内容块
func PDFPart(base64Data string) ContentPart {
	return ContentPart{Type: ContentPDF, Base64: "data:application/pdf;base64," + base64Data}
}

// ParseBase64DataURL 解析 data:xxx;base64,... 格式
// 返回 MIME 类型和 base64 数据（不含前缀）
func ParseBase64DataURL(dataURL string) (mimeType string, data string) {
	if len(dataURL) < 5 || dataURL[:5] != "data:" {
		return "", ""
	}

	// 找到 base64, 的位置
	commaIdx := -1
	for i := 5; i < len(dataURL); i++ {
		if dataURL[i] == ',' {
			commaIdx = i
			break
		}
	}
	if commaIdx == -1 {
		return "", ""
	}

	// 解析 MIME 类型
	header := dataURL[5:commaIdx] // 去掉 "data:"
	semicolonIdx := -1
	for i := 0; i < len(header); i++ {
		if header[i] == ';' {
			semicolonIdx = i
			break
		}
	}
	if semicolonIdx != -1 {
		mimeType = header[:semicolonIdx]
	} else {
		mimeType = header
	}

	return mimeType, dataURL[commaIdx+1:]
}
