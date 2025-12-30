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

type OpenAIProvider struct {
	client     *openai.Client
	httpClient *http.Client
	model      string
	baseURL    string
	apiKey     string
}

func NewOpenAIProvider(apiKey, baseURL, model string) *OpenAIProvider {
	if model == "" {
		model = openai.ChatModelGPT4o
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{},
	}

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

	return &OpenAIProvider{
		client:     &client,
		httpClient: httpClient,
		model:      model,
		baseURL:    baseURL,
		apiKey:     apiKey,
	}
}

func (p *OpenAIProvider) GenerateContentStream(ctx context.Context, history []openai.ChatCompletionMessageParamUnion, onToken func(string)) (string, error) {
	stream := p.client.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
		Model:    p.model,
		Messages: history,
	})

	defer stream.Close()

	var fullResponse strings.Builder

	for stream.Next() {
		evt := stream.Current()

		if len(evt.Choices) > 0 {
			content := evt.Choices[0].Delta.Content
			if content != "" {
				fullResponse.WriteString(content)
				if onToken != nil {
					onToken(content)
				}
			}
		}
	}

	if err := stream.Err(); err != nil {
		errStr := err.Error()

		startIndex := strings.Index(errStr, "{")
		if startIndex == -1 {
			return "", fmt.Errorf("未知错误: %s", errStr)
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
			return "", fmt.Errorf("解析错误: %s", response.Message)
		}

		return "", fmt.Errorf("%s", string(finalJsonBytes))
	}

	return fullResponse.String(), nil
}

func (p *OpenAIProvider) GetBalance(ctx context.Context) (float64, error) {
	_, err := p.client.Models.List(ctx)
	if err != nil {
		return 0, err
	}
	return 0, nil
}

func (p *OpenAIProvider) Validate(ctx context.Context) error {
	_, err := p.client.Models.List(ctx)
	return err
}

func (p *OpenAIProvider) GetModels(ctx context.Context) ([]string, error) {
	resp, err := p.client.Models.List(ctx)
	if err != nil {
		return nil, err
	}
	var models []string
	for _, m := range resp.Data {
		models = append(models, m.ID)
	}
	return models, nil
}

func (p *OpenAIProvider) ParseResume(ctx context.Context, resumeBase64 string) (string, error) {
	messages := []openai.ChatCompletionMessageParamUnion{}
	messages = append(messages, openai.UserMessage(`# Role 你是一个**通用型简历重构与解析引擎**。你的任务是将输入的简历内容（无论行业是互联网、销售、财务还是行政），提取并严格按照我指定的【通用美观模板】转换为 Raw Markdown 格式。 # 🚨 STRICT OUTPUT PROTOCOL (绝对输出协议) 1. **纯文本输出**：输出**必须**是原始 Markdown 文本。**严禁**使用 markdown 代码块包裹。 2. **内容自适应**：根据候选人的行业调整关键词。例如： - 对于程序员，提取"技术栈"； - 对于销售，提取"关键客户/业绩"； - 对于行政，提取"办公技能/组织能力"。 3. **空值处理**：如果简历中缺少某项信息（如个人网站），直接不显示该行。 4. **排版强制**：严格保留模板中的 Emoji 图标和引用块格式。 # 💅 Universal Visual Template (通用美观模板) # {{姓名}} > 💼 **{{求职意向/当前职位}}** > > 📱 {{电话}}  |  📧 {{邮箱}}  |  📍 {{所在城市}} > 🔗 [作品集/LinkedIn/个人主页]({{链接地址}}) *(如有才显示)* --- ## ⚡ 专业技能 *(请根据职业属性归类，以下仅为示例，请灵活调整 Key)* - **核心竞争力**: {{如：大客户销售 / 财务审计 / Java开发 / 团队管理}} - **软件/工具**: {{如：SAP / Excel高级 / Photoshop / Docker}} - **证书/语言**: {{如：CPA注册会计师 / 英语六级 / PMP}} ## 🏢 工作经历 ### **{{公司名称}}** **{{职位名称}}** | *{{开始时间}} - {{结束时间}}* > {{一句话概括核心职责。如：负责大区销售团队管理，或：负责公司年度审计工作。}} - 🔸 **{{核心业绩/产出1}}**: {{详细描述，尽可能包含数据，如：销售额增长 20% / 节省成本 50万}} - 🔸 **{{核心业绩/产出2}}**: {{详细描述}} - 🔸 **{{核心业绩/产出3}}**: {{详细描述}} *(如果有更多公司，重复上面的格式)* ## 🏆 项目与重点业绩 *(如果是技术人员写项目；如果是销售写大客户案例；如果是应届生写校园活动)* ### 🔹 {{项目/案例名称}} *{{项目/案例一句话简介}}* - **扮演角色**: {{如：项目负责人 / 核心执行者}} - **背景/挑战**: {{简述面临的问题}} - **我的行动**: - {{具体动作 1}} - {{具体动作 2}} - **最终结果**: {{量化的结果}} ## 🎓 教育经历 - **{{学校名称}}** | {{专业}} | {{学历}} | *{{时间段}}* --- *Generated by AI Resume Assistant* # Input Data 简历内容见附件。`))
	messages = append(messages, openai.UserMessage([]openai.ChatCompletionContentPartUnionParam{
		openai.ImageContentPart(openai.ChatCompletionContentPartImageImageURLParam{
			URL: "data:application/pdf;base64," + resumeBase64,
		}),
	}))
	result, err := p.GenerateContentStream(ctx, messages, nil)
	if err != nil {
		return "", err
	}
	return result, nil
}
