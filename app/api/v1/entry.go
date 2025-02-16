package v1

import "go-server-base/app/service"

type ApiGroup struct {
	BaseApi
}

var ApiGroupApp = new(ApiGroup)

var (
	ragService = service.NewIRagService()
)
