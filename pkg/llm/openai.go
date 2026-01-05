package llm

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	openai "github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// OpenAIAdapter OpenAI 适配器
type OpenAIAdapter struct {
	client     *openai.Client
	httpClient *http.Client
	model      string
	baseURL    string
	apiKey     string
}

// NewOpenAIAdapter 创建 OpenAI 适配器
func NewOpenAIAdapter(apiKey, baseURL, model string) *OpenAIAdapter {
	if model == "" {
		model = openai.ChatModelGPT4o
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{},
	}
	// u, _ := url.Parse("http://127.0.0.1:8888")
	// transport.Proxy = http.ProxyURL(u)
	httpClient := &http.Client{
		Transport: transport,
	}

	opts := []option.RequestOption{
		option.WithAPIKey(apiKey),
		option.WithHTTPClient(httpClient),
	}

	if baseURL != "" {
		opts = append(opts, option.WithBaseURL(baseURL))
	}

	client := openai.NewClient(opts...)

	return &OpenAIAdapter{
		client:     &client,
		httpClient: httpClient,
		model:      model,
		baseURL:    baseURL,
		apiKey:     apiKey,
	}
}

// ==================== 类型转换方法 ====================

// toOpenAIMessages 将统一格式转换为 OpenAI SDK 格式
func (a *OpenAIAdapter) toOpenAIMessages(messages []Message) []openai.ChatCompletionMessageParamUnion {
	result := make([]openai.ChatCompletionMessageParamUnion, 0, len(messages))

	for _, msg := range messages {
		switch msg.Role {
		case RoleSystem:
			result = append(result, openai.SystemMessage(msg.Content))

		case RoleUser:
			if len(msg.Parts) > 0 {
				result = append(result, openai.UserMessage(a.toOpenAIParts(msg.Parts)))
			} else {
				result = append(result, openai.UserMessage(msg.Content))
			}

		case RoleAssistant:
			result = append(result, openai.AssistantMessage(msg.Content))
		}
	}

	return result
}

// toOpenAIParts 将 ContentPart 转换为 OpenAI 格式
func (a *OpenAIAdapter) toOpenAIParts(parts []ContentPart) []openai.ChatCompletionContentPartUnionParam {
	result := make([]openai.ChatCompletionContentPartUnionParam, 0, len(parts))

	for _, part := range parts {
		switch part.Type {
		case ContentText:
			result = append(result, openai.TextContentPart(part.Text))
		case ContentImage, ContentPDF:
			result = append(result, openai.ImageContentPart(openai.ChatCompletionContentPartImageImageURLParam{
				URL: part.Base64,
			}))
		}
	}

	return result
}

// ==================== Provider 接口实现 ====================

// GenerateContentStream 流式生成内容
func (a *OpenAIAdapter) GenerateContentStream(ctx context.Context, messages []Message, onChunk StreamCallback) (Message, error) {
	openaiMessages := a.toOpenAIMessages(messages)

	stream := a.client.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
		Model:    a.model,
		Messages: openaiMessages,
	})

	defer stream.Close()

	var fullContent strings.Builder
	var fullThinking strings.Builder

	for stream.Next() {
		evt := stream.Current()

		if len(evt.Choices) > 0 {
			delta := evt.Choices[0].Delta
			content := delta.Content

			if content != "" {
				fullContent.WriteString(content)

				if onChunk != nil {
					onChunk(StreamChunk{
						Type:    ChunkContent,
						Content: content,
					})
				}
			}
		}
	}

	if err := stream.Err(); err != nil {
		return Message{}, a.parseError(err)
	}

	return Message{
		Role:     RoleAssistant,
		Content:  fullContent.String(),
		Thinking: fullThinking.String(),
	}, nil
}

// parseError 解析错误信息
func (a *OpenAIAdapter) parseError(err error) error {
	errStr := err.Error()

	startIndex := strings.Index(errStr, "{")
	if startIndex == -1 {
		return fmt.Errorf("未知错误: %s", errStr)
	}

	jsonPart := errStr[startIndex:]
	var response struct {
		StatusCode int    `json:"statusCode"`
		Code       string `json:"code"`
		Message    string `json:"message"`
		Type       string `json:"type"`
	}

	_ = json.Unmarshal([]byte(jsonPart), &response)
	headerPart := errStr[:startIndex]

	lastColon := strings.LastIndex(headerPart, ":")
	if lastColon != -1 {
		statusPart := headerPart[lastColon+1:]
		if _, scanErr := fmt.Sscanf(statusPart, "%d", &response.StatusCode); scanErr != nil {
			response.StatusCode = 500
		}
	} else {
		response.StatusCode = -1
	}

	finalJsonBytes, marshalErr := json.Marshal(response)
	if marshalErr != nil {
		return fmt.Errorf("解析错误: %s", response.Message)
	}

	return fmt.Errorf("%s", string(finalJsonBytes))
}

// TestChat 测试连通性
func (a *OpenAIAdapter) TestChat(ctx context.Context) error {
	_, err := a.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: a.model,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage("hi"),
		},
		MaxTokens: openai.Int(1),
	})
	return err
}

// GetModels 获取模型列表
func (a *OpenAIAdapter) GetModels(ctx context.Context) ([]string, error) {
	resp, err := a.client.Models.List(ctx)
	if err != nil {
		return nil, err
	}
	var models []string
	for _, m := range resp.Data {
		models = append(models, m.ID)
	}
	return models, nil
}

// ParseResume 解析简历为 Markdown
func (a *OpenAIAdapter) ParseResume(ctx context.Context, resumeBase64 string) (string, error) {
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
