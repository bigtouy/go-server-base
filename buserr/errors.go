package buserr

import (
	"go-server-base/i18n"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type BusinessError struct {
	Msg    string
	Detail interface{}
	Map    map[string]interface{}
	Err    error
}

func (e BusinessError) Error(ctx *gin.Context) string {
	content := ""
	if e.Detail != nil {
		content = i18n.GetErrMsg(ctx, e.Msg, map[string]interface{}{"detail": e.Detail})
	} else if e.Map != nil {
		content = i18n.GetErrMsg(ctx, e.Msg, e.Map)
	} else {
		content = i18n.GetErrMsg(ctx, e.Msg, nil)
	}
	if content == "" {
		if e.Err != nil {
			return e.Err.Error()
		}
		return errors.New(e.Msg).Error()
	}
	return content
}

func New(Key string) BusinessError {
	return BusinessError{
		Msg:    Key,
		Detail: nil,
		Err:    nil,
	}
}

func WithDetail(Key string, detail interface{}, err error) BusinessError {
	return BusinessError{
		Msg:    Key,
		Detail: detail,
		Err:    err,
	}
}

func WithErr(Key string, err error) BusinessError {
	return BusinessError{
		Msg:    Key,
		Detail: "",
		Err:    err,
	}
}

func WithMap(Key string, maps map[string]interface{}, err error) BusinessError {
	return BusinessError{
		Msg: Key,
		Map: maps,
		Err: err,
	}
}
