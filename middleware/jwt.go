package middleware

import (
	"go-server-base/app/api/v1/helper"
	"go-server-base/constant"
	jwtUtils "go-server-base/utils/jwt"

	"github.com/gin-gonic/gin"
)

func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get(constant.JWTHeaderName)
		if token == "" {
			c.Next()
			return
		}
		j := jwtUtils.NewJWT()
		claims, err := j.ParseToken(token)
		if err != nil {
			helper.ErrorWithDetail(c, constant.CodeErrUnauthorized, constant.ErrTypeInternalServer, err)
			return
		}
		c.Set("claims", claims)
		c.Set("authMethod", constant.AuthMethodJWT)
		c.Next()
	}
}
