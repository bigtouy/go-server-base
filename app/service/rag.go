package service

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dianjin20240628 "github.com/alibabacloud-go/dianjin-20240628/client"
	openapiutil "github.com/alibabacloud-go/openapi-util/service"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/gin-gonic/gin"
	"go-server-base/global"
	"net/http"
)

type RagService struct {
	BaseApi  string            `json:"base_api"`
	Endpoint string            `json:"endpoint"`
	Client   http.Client       `json:"client"`
	RagMap   map[string]string `json:"rag_map"`
}

type RagResponse struct {
	Headers    map[string]*string     `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                 `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       map[string]interface{} `json:"body,omitempty" xml:"body,omitempty"`
}

type ErrResponse struct {
	StatusCode *int32                 `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Data       map[string]interface{} `json:"data,omitempty" xml:"data,omitempty"`
}

type IRagService interface {
	RequestRagApi(c *gin.Context) (*RagResponse, error)
}

var ragApiMap = map[string]string{
	"/api/app/config":                                "GetAppConfig",
	"/api/library/create":                            "CreateLibrary",
	"/api/library/document/getParseResult":           "GetParseResult",
	"/api/library/list":                              "GetLibraryList",
	"/api/library/get":                               "GetLibrary",
	"/api/library/document/upload":                   "UploadDocument",
	"/api/library/document/url":                      "GetDocumentUrl",
	"/api/library/document/preview":                  "PreviewDocument",
	"/api/library/filterDocument":                    "GetFilterDocumentList",
	"/api/library/listDocument":                      "GetDocumentList",
	"/api/library/document/delete":                   "DeleteDocument",
	"/api/library/document/updateDocument":           "UpdateDocument",
	"/api/library/getDocumentChunk":                  "GetDocumentChunkList",
	"/api/library/document/createPredefinedDocument": "CreatePredefinedDocument",
	"/api/library/recallDocument":                    "RecallDocument",
	"/api/library/document/reIndex":                  "ReIndex",
	"/api/library/update":                            "UpdateLibrary",
	"/api/library/delete":                            "DeleteLibrary",
	"/api/run/library/chat/generation":               "RunLibraryChatGeneration",
	"/api/history/list":                              "GetHistoryListByBizType",
	"/api/plugin/invoke":                             "InvokePlugin",
}

func NewIRagService() IRagService {
	return &RagService{
		BaseApi:  "https://dianjin.cn-beijing.aliyuncs.com",
		Endpoint: "dianjin.cn-beijing.aliyuncs.com",
		Client:   http.Client{},
		RagMap:   ragApiMap,
	}
}

func (r RagService) RequestRagApi(c *gin.Context) (*RagResponse, error) {
	config := &openapi.Config{
		// 必填，请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_ID。
		AccessKeyId: tea.String(global.CONF.Oss.AccessKeyId),
		// 必填，请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_SECRET。
		AccessKeySecret: tea.String(global.CONF.Oss.AccessKeySecret),
	}
	// Endpoint 请参考 https://api.aliyun.com/product/DianJin
	config.Endpoint = tea.String(r.Endpoint)
	client, err := dianjin20240628.NewClient(config)
	if err != nil {
		return nil, tea.NewSDKError(map[string]interface{}{})
	}
	headers := make(map[string]*string)

	workspaceId := tea.String(global.CONF.Qwen.WorkspaceId)
	query := c.Request.URL.Query()
	global.LOG.Info("query", query)

	currentQuery := make(map[string]*string)

	for k, v := range query {
		currentQuery[k] = &v[0]
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    c.Request.Body,
		Query:   currentQuery,
	}

	apiPath := c.Param("action")
	action := r.RagMap[c.Param("action")]

	params := &openapi.Params{
		Action:      tea.String(action),
		Version:     tea.String("2024-06-28"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/" + tea.StringValue(openapiutil.GetEncodeParam(workspaceId)) + apiPath),
		Method:      tea.String(c.Request.Method),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}

	runtime := &util.RuntimeOptions{}
	_result := &RagResponse{}
	response, err := client.DoRequest(params, req, runtime)

	if err != nil {
		return nil, err
	}

	err = tea.Convert(response, &_result)

	if err != nil {
		return nil, err
	}
	return _result, err
}
