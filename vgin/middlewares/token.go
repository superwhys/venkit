package middlewares

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/superwhys/venkit/authutils"
	"github.com/superwhys/venkit/lg"
)

const tokenKey = "auth-token"

func TokenManagerMiddleware(tokenTmpl authutils.Token, tokenManager *authutils.TokenManager) gin.HandlerFunc {
	t := reflect.TypeOf(tokenTmpl)
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		lg.Fatal("NewTaskQueue: typeObj should be ptr to struct")
	}

	t = t.Elem()

	return func(c *gin.Context) {
		headerToken := c.GetHeader(AuthHeaderKey)

		newToken := reflect.New(t).Interface().(authutils.Token)
		newToken.SetKey(headerToken)

		if err := tokenManager.Read(newToken); err != nil {
			lg.Errorf("token read error: %v", err)
			vgin.AbortWithError(c, http.StatusUnauthorized, "认证失败，请求需要token")
			return
		}

		c.Set(tokenKey, newToken)
	}
}

func GetToken(c *gin.Context) authutils.Token {
	val, exists := c.Get(tokenKey)
	if !exists {
		return nil
	}

	return val.(authutils.Token)
}
