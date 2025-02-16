package v1

import (
	"encoding/json"
	"errors"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/gin-gonic/gin"
	"go-server-base/global"
	"reflect"
)

type BaseApi struct{}

type ErrResponse struct {
	StatusCode *int32 `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Data       string `json:"data,omitempty" xml:"data,omitempty"`
}

type ResData struct {
	Cost       *interface{} `json:"cost"`
	Data       *interface{} `json:"data"`
	DataType   *interface{} `json:"dataType"`
	ErrCode    string       `json:"errCode"`
	Message    string       `json:"message"`
	RequestId  string       `json:"requestId"`
	StatusCode int          `json:"statusCode"`
	Success    bool         `json:"success"`
	Time       *interface{} `json:"time"`
}

func newResData() ResData {
	return ResData{
		Cost:       nil,
		Data:       nil,
		DataType:   nil,
		ErrCode:    "400",
		Message:    "success",
		RequestId:  "",
		StatusCode: 400,
		Success:    false,
		Time:       nil,
	}
}

func typeToCustomError(err error) (*tea.SDKError, bool) {
	var customErr *tea.SDKError
	ok := errors.As(err, &customErr)
	return customErr, ok
}

func (b *BaseApi) Rag(c *gin.Context) {

	response, err := ragService.RequestRagApi(c)

	if err != nil {
		global.LOG.Error(err)
		resData := newResData()
		if reflect.ValueOf(err).Type().String() == "*tea.SDKError" {
			customErr, ok := typeToCustomError(err)

			if ok == false {
				resData.Message = err.Error()
				c.JSON(200, resData)
			}

			dataMap := make(map[string]interface{})
			err := json.Unmarshal([]byte(*customErr.Data), &dataMap)
			if err != nil {
				resData.Message = err.Error()
				c.JSON(200, resData)
				return
			}
			c.JSON(200, dataMap)
			return
		}

		resData.Message = err.Error()
		c.JSON(200, resData)
		return
	}
	global.LOG.Info(response)

	for k, v := range response.Headers {
		c.Header(k, *v)
	}

	c.JSON(int(*response.StatusCode), response.Body)
}
