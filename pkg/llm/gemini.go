package llm

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"google.golang.org/genai"
)

// GeminiAdapter Gemini 适配器
type GeminiAdapter struct {
	client *genai.Client
	model  string
	apiKey string
}

// NewGeminiAdapter 创建 Gemini 适配器
func NewGeminiAdapter(apiKey, model string) (*GeminiAdapter, error) {
	if model == "" {
		model = "gemini-2.0-flash"
	}

	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	return &GeminiAdapter{
		client: client,
		model:  model,
		apiKey: apiKey,
	}, nil
}

// ==================== 类型转换方法 ====================

// toGeminiContents 将统一格式转换为 Gemini SDK 格式
func (a *GeminiAdapter) toGeminiContents(messages []Message) ([]*genai.Content, string) {
	var contents []*genai.Content
	var systemInstruction string

	for _, msg := range messages {
		switch msg.Role {
		case RoleSystem:
			systemInstruction = msg.Content

		case RoleUser:
			parts := a.toGeminiParts(msg)
			contents = append(contents, &genai.Content{
				Role:  "user",
				Parts: parts,
			})

		case RoleAssistant:
			contents = append(contents, &genai.Content{
				Role:  "model",
				Parts: []*genai.Part{{Text: msg.Content}},
			})
		}
	}

	return contents, systemInstruction
}

// toGeminiParts 将 Message 转换为 Gemini Parts
func (a *GeminiAdapter) toGeminiParts(msg Message) []*genai.Part {
	if len(msg.Parts) == 0 {
		return []*genai.Part{{Text: msg.Content}}
	}

	var parts []*genai.Part
	for _, p := range msg.Parts {
		switch p.Type {
		case ContentText:
			parts = append(parts, &genai.Part{Text: p.Text})

		case ContentImage:
			// 解析 base64 图片
			mimeType, data := parseBase64DataURL(p.Base64)
			if data != nil {
				parts = append(parts, &genai.Part{
					InlineData: &genai.Blob{
						MIMEType: mimeType,
						Data:     data,
					},
				})
			}

		case ContentPDF:
			// 解析 base64 PDF
			_, data := parseBase64DataURL(p.Base64)
			if data != nil {
				parts = append(parts, &genai.Part{
					InlineData: &genai.Blob{
						MIMEType: "application/pdf",
						Data:     data,
					},
				})
			}
		}
	}

	return parts
}

// parseBase64DataURL 解析 data:xxx;base64,... 格式
func parseBase64DataURL(dataURL string) (mimeType string, data []byte) {
	// 格式: data:image/png;base64,xxxxxx
	if !strings.HasPrefix(dataURL, "data:") {
		return "", nil
	}

	// 找到 base64, 的位置
	commaIdx := strings.Index(dataURL, ",")
	if commaIdx == -1 {
		return "", nil
	}

	// 解析 MIME 类型
	header := dataURL[5:commaIdx] // 去掉 "data:"
	semicolonIdx := strings.Index(header, ";")
	if semicolonIdx != -1 {
		mimeType = header[:semicolonIdx]
	} else {
		mimeType = header
	}

	// 解码 base64
	base64Data := dataURL[commaIdx+1:]
	decoded, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return "", nil
	}

	return mimeType, decoded
}

// ==================== Provider 接口实现 ====================

// GenerateContentStream 流式生成内容
func (a *GeminiAdapter) GenerateContentStream(ctx context.Context, messages []Message, onChunk StreamCallback) (Message, error) {
	contents, systemInstruction := a.toGeminiContents(messages)

	config := &genai.GenerateContentConfig{
		ThinkingConfig: &genai.ThinkingConfig{
			IncludeThoughts: true,
		},
	}
	if systemInstruction != "" {
		config.SystemInstruction = &genai.Content{
			Parts: []*genai.Part{{Text: systemInstruction}},
		}
	}

	// 使用流式生成
	var fullContent strings.Builder
	var fullThinking strings.Builder

	for resp := range a.client.Models.GenerateContentStream(ctx, a.model, contents, config) {
		if resp == nil {
			continue
		}

		// 遍历所有候选响应
		for _, candidate := range resp.Candidates {
			if candidate.Content == nil {
				continue
			}

			for _, part := range candidate.Content.Parts {
				if part == nil {
					continue
				}

				// 检查是否是思维链（Gemini 2.0 Flash Thinking 模型）
				if part.Thought {
					fullThinking.WriteString(part.Text)
					if onChunk != nil {
						onChunk(StreamChunk{
							Type:    ChunkThinking,
							Content: part.Text,
						})
					}
				} else if part.Text != "" {
					fullContent.WriteString(part.Text)
					if onChunk != nil {
						onChunk(StreamChunk{
							Type:    ChunkContent,
							Content: part.Text,
						})
					}
				}
			}
		}
	}

	// 返回最终结果
	return Message{
		Role:     RoleAssistant,
		Content:  fullContent.String(),
		Thinking: fullThinking.String(),
	}, nil
}

// TestChat 测试连通性
func (a *GeminiAdapter) TestChat(ctx context.Context) error {
	contents := []*genai.Content{
		{
			Role:  "user",
			Parts: []*genai.Part{{Text: "hi"}},
		},
	}

	config := &genai.GenerateContentConfig{
		MaxOutputTokens: 1,
	}

	_, err := a.client.Models.GenerateContent(ctx, a.model, contents, config)
	return err
}

// GetModels 获取模型列表
func (a *GeminiAdapter) GetModels(ctx context.Context) ([]string, error) {
	page, err := a.client.Models.List(ctx, &genai.ListModelsConfig{})
	if err != nil {
		return nil, err
	}

	var models []string
	for _, model := range page.Items {
		if model != nil {
			models = append(models, model.Name)
		}
	}

	return models, nil
}

// ParseResume 解析简历为 Markdown
func (a *GeminiAdapter) ParseResume(ctx context.Context, resumeBase64 string) (string, error) {
	messages := []Message{
		NewUserMessage(`# Role 你是一个**通用型简历重构与解析引擎**。你的任务是将输入的简历内容（无论行业是互联网、销售、财务还是行政），提取并严格按照我指定的【通用美观模板】转换为 Raw Markdown 格式。 # 🚨 STRICT OUTPUT PROTOCOL (绝对输出协议) 1. **纯文本输出**：输出**必须**是原始 Markdown 文本。**严禁**使用 markdown 代码块包裹。 2. **内容自适应**：根据候选人的行业调整关键词。例如： - 对于程序员，提取"技术栈"； - 对于销售，提取"关键客户/业绩"； - 对于行政，提取"办公技能/组织能力"。 3. **空值处理**：如果简历中缺少某项信息（如个人网站），直接不显示该行。 4. **排版强制**：严格保留模板中的 Emoji 图标和引用块格式。 # 💅 Universal Visual Template (通用美观模板) # {{姓名}} > 💼 **{{求职意向/当前职位}}** > > 📱 {{电话}}  |  📧 {{邮箱}}  |  📍 {{所在城市}} > 🔗 [作品集/LinkedIn/个人主页]({{链接地址}}) *(如有才显示)* --- ## ⚡ 专业技能 *(请根据职业属性归类，以下仅为示例，请灵活调整 Key)* - **核心竞争力**: {{如：大客户销售 / 财务审计 / Java开发 / 团队管理}} - **软件/工具**: {{如：SAP / Excel高级 / Photoshop / Docker}} - **证书/语言**: {{如：CPA注册会计师 / 英语六级 / PMP}} ## 🏢 工作经历 ### **{{公司名称}}** **{{职位名称}}** | *{{开始时间}} - {{结束时间}}* > {{一句话概括核心职责。如：负责大区销售团队管理，或：负责公司年度审计工作。}} - 🔸 **{{核心业绩/产出1}}**: {{详细描述，尽可能包含数据，如：销售额增长 20% / 节省成本 50万}} - 🔸 **{{核心业绩/产出2}}**: {{详细描述}} - 🔸 **{{核心业绩/产出3}}**: {{详细描述}} *(如果有更多公司，重复上面的格式)* ## 🏆 项目与重点业绩 *(如果是技术人员写项目；如果是销售写大客户案例；如果是应届生写校园活动)* ### 🔹 {{项目/案例名称}} *{{项目/案例一句话简介}}* - **扮演角色**: {{如：项目负责人 / 核心执行者}} - **背景/挑战**: {{简述面临的问题}} - **我的行动**: - {{具体动作 1}} - {{具体动作 2}} - **最终结果**: {{量化的结果}} ## 🎓 教育经历 - **{{学校名称}}** | {{专业}} | {{学历}} | *{{时间段}}* --- *Generated by AI Resume Assistant* # Input Data 简历内容见附件。`),
		NewMultiPartMessage(RoleUser, []ContentPart{
			PDFPart(resumeBase64),
		}),
	}

	result, err := a.GenerateContentStream(ctx, messages, nil)
	if err != nil {
		return "", err
	}
	return result.Content, nil
}
