package middleware

import (
	"strconv"
	"time"

	"go-server-base/app/api/v1/helper"
	"go-server-base/app/repo"
	"go-server-base/constant"

	"github.com/gin-gonic/gin"
)

func PasswordExpired() gin.HandlerFunc {
	return func(c *gin.Context) {
		settingRepo := repo.NewISettingRepo()
		setting, err := settingRepo.Get(settingRepo.WithByKey("ExpirationDays"))
		if err != nil {
			helper.ErrorWithDetail(c, constant.CodePasswordExpired, constant.ErrTypePasswordExpired, err)
			return
		}
		expiredDays, _ := strconv.Atoi(setting.Value)
		if expiredDays == 0 {
			c.Next()
			return
		}

		extime, err := settingRepo.Get(settingRepo.WithByKey("ExpirationTime"))
		if err != nil {
			helper.ErrorWithDetail(c, constant.CodePasswordExpired, constant.ErrTypePasswordExpired, err)
			return
		}
		expiredTime, err := time.Parse("2006-01-02 15:04:05", extime.Value)
		if err != nil {
			helper.ErrorWithDetail(c, constant.CodePasswordExpired, constant.ErrTypePasswordExpired, err)
			return
		}
		if time.Now().After(expiredTime) {
			helper.ErrorWithDetail(c, constant.CodePasswordExpired, constant.ErrTypePasswordExpired, nil)
			return
		}
		c.Next()
	}
}
