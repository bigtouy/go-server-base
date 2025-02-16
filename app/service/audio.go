package service

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/gin-gonic/gin"
	"go-server-base/global"
	"time"
)

type AreaParams struct {
	EndpointName      string
	Product           string
	Domain            string
	ApiVersion        string
	PostRequestAction string
	GetRequestAction  string
}

type RequestParams struct {
	AppKey      string
	Version     string
	EnableWords string
}

type AudioService struct {
	Client        *sdk.Client
	AreaParams    AreaParams
	RequestParams RequestParams
}

type IAudioService interface {
	GetAudio(c *gin.Context, fileLink string)
}

const (
	success  = "SUCCESS"
	running  = "RUNNING"
	queueing = "QUEUEING"
)

func NewIAudioService() IAudioService {
	credential := credentials.AccessKeyCredential{
		AccessKeyId:     global.CONF.Oss.AccessKeyId,
		AccessKeySecret: global.CONF.Oss.AccessKeySecret,
	}
	client, err := sdk.NewClientWithOptions(global.CONF.Qwen.Region, &sdk.Config{}, credential)
	if err != nil {
		panic(err)
	}

	return &AudioService{
		Client: client,
		AreaParams: AreaParams{
			EndpointName:      "cn-shanghai",
			Product:           "nls-filetrans",
			Domain:            "filetrans.cn-shanghai.aliyuncs.com",
			ApiVersion:        "2018-08-17",
			PostRequestAction: "SubmitTask",
			GetRequestAction:  "GetTaskResult",
		},
		RequestParams: RequestParams{
			AppKey:      global.CONF.Qwen.NlsAppKey,
			Version:     "4.0",
			EnableWords: "false",
		},
	}
}

func (u AudioService) GetAudio(c *gin.Context, fileLink string) {
	postRequest := requests.NewCommonRequest()
	postRequest.Domain = u.AreaParams.Domain
	postRequest.Version = u.AreaParams.ApiVersion
	postRequest.Product = u.AreaParams.Product
	postRequest.ApiName = u.AreaParams.PostRequestAction
	postRequest.Method = "POST"
	mapTask := make(map[string]string)
	mapTask["appkey"] = u.RequestParams.AppKey
	mapTask["file_link"] = fileLink
	// 新接入请使用4.0版本，已接入（默认2.0）如需维持现状，请注释掉该参数设置。
	mapTask["version"] = u.RequestParams.Version
	// 设置是否输出词信息，默认为false。开启时需要设置version为4.0。
	mapTask["enable_words"] = u.RequestParams.EnableWords
	task, err := json.Marshal(mapTask)
	if err != nil {
		panic(err)
	}
	postRequest.FormParams["Task"] = string(task)
	postResponse, err := u.Client.ProcessCommonRequest(postRequest)
	if err != nil {
		panic(err)
	}
	postResponseContent := postResponse.GetHttpContentString()
	fmt.Println(postResponseContent)
	if postResponse.GetHttpStatus() != 200 {
		fmt.Println("录音文件识别请求失败，Http错误码: ", postResponse.GetHttpStatus())
		return
	}
	var postMapResult map[string]interface{}
	err = json.Unmarshal([]byte(postResponseContent), &postMapResult)
	if err != nil {
		panic(err)
	}
	var taskId = ""
	var statusText = ""
	statusText = postMapResult["StatusText"].(string)
	if statusText == success {
		fmt.Println("录音文件识别请求成功响应!")
		taskId = postMapResult["TaskId"].(string)
	} else {
		fmt.Println("录音文件识别请求失败!")
		return
	}
	getRequest := requests.NewCommonRequest()
	getRequest.Domain = u.AreaParams.Domain
	getRequest.Version = u.AreaParams.ApiVersion
	getRequest.Product = u.AreaParams.Product
	getRequest.ApiName = u.AreaParams.PostRequestAction
	getRequest.Method = "GET"
	getRequest.QueryParams["TaskId"] = taskId
	statusText = ""
	for {
		getResponse, err := u.Client.ProcessCommonRequest(getRequest)
		if err != nil {
			panic(err)
		}
		getResponseContent := getResponse.GetHttpContentString()
		fmt.Println("识别查询结果：", getResponseContent)
		if getResponse.GetHttpStatus() != 200 {
			fmt.Println("识别结果查询请求失败，Http错误码：", getResponse.GetHttpStatus())
			break
		}
		var getMapResult map[string]interface{}
		err = json.Unmarshal([]byte(getResponseContent), &getMapResult)
		if err != nil {
			panic(err)
		}
		statusText = getMapResult["StatusText"].(string)
		if statusText == running || statusText == queueing {
			time.Sleep(10 * time.Second)
		} else {
			break
		}

	}
	if statusText == success {
		fmt.Println("录音文件识别成功！")
	} else {
		fmt.Println("录音文件识别失败！")
	}
}
