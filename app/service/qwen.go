package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-server-base/global"
	"io"
	"net/http"
)

type QwenService struct {
	BaseApi string `json:"base_api"`
	Method  string `json:"method"`
	Model   string `json:"model"`
	ApiKey  string `json:"api_key"`
}

type IQwenService interface {
	VideoRecogByUrl(c *gin.Context, url string) (string, error)
	VideoRecogByImages(c *gin.Context, images []string) (string, error)
}

func NewIQwenService() IQwenService {
	return &QwenService{
		BaseApi: "https://dashscope.aliyuncs.com/api/v1/services/aigc/multimodal-generation/generation",
		Method:  "post",
		Model:   "qwen-vl-max",
		ApiKey:  global.CONF.Qwen.ApiKey,
	}
}

func (q QwenService) VideoRecogByUrl(c *gin.Context, url string) (string, error) {
	baseApi := q.BaseApi
	method := q.Method
	model := q.Model
	apiKey := q.ApiKey

	prompt := `
# ROLE 
你是一个视频理解专家
# WORKFLOW
1. 先完整理解视频内容，然后描述视频内容，最后采用总分结构输出
2. 严格按照总分结构，第一句提供一个简短的内容总结。 第二句话开始详细描述视频的具体内容，越详细越好。
# REMEMBER
描述视频内容越详细越好，禁止repeat，禁止出现重复内容
`

	// 构建请求体
	requestBody := map[string]interface{}{
		"model": model,
		"input": map[string]interface{}{
			"messages": []map[string]interface{}{
				{
					"role": "user",
					"content": []map[string]interface{}{
						{
							"video": url,
						},
						{
							"text": prompt,
						},
					},
				},
			},
		},
	}

	jsonBytes, err := json.Marshal(requestBody)

	if err != nil {
		global.LOG.Error("Error marshalling request body:", err)
		return "", nil
	}

	payload := bytes.NewBuffer(jsonBytes)

	client := &http.Client{}
	req, err := http.NewRequest(method, baseApi, payload)

	if err != nil {
		fmt.Println(err)
		return "", nil
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Connection", "keep-alive")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", nil
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "", nil
	}
	fmt.Println(string(body))

	return "", nil
}

func (q QwenService) VideoRecogByImages(c *gin.Context, images []string) (string, error) {
	baseApi := q.BaseApi
	method := q.Method
	model := q.Model
	apiKey := q.ApiKey
	prompt := `
# ROLE 
你是一个视频理解专家
# WORKFLOW
1. 先完整理解视频内容，然后描述视频内容，最后采用总分结构输出
2. 严格按照总分结构，第一句提供一个简短的内容总结。 第二句话开始详细描述视频的具体内容，越详细越好。
# REMEMBER
描述视频内容越详细越好，禁止repeat，禁止出现重复内容
`
	// 构建请求体
	requestBody := map[string]interface{}{
		"model": model,
		"input": map[string]interface{}{
			"messages": []map[string]interface{}{
				{
					"role": "user",
					"content": []map[string]interface{}{
						{
							"video": images,
						},
						{
							"text": prompt,
						},
					},
				},
			},
		},
	}

	jsonBytes, err := json.Marshal(requestBody)

	if err != nil {
		global.LOG.Error("Error marshalling request body:", err)
		return "", nil
	}

	payload := bytes.NewBuffer(jsonBytes)

	client := &http.Client{}
	req, err := http.NewRequest(method, baseApi, payload)

	if err != nil {
		fmt.Println(err)
		return "", nil
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Connection", "keep-alive")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", nil
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "", nil
	}
	fmt.Println(string(body))

	return "", nil
}
