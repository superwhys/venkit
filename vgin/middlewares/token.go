package middlewares

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/superwhys/venkit/lg"
	"github.com/superwhys/venkit/vauth"
	"github.com/superwhys/venkit/vgin"
)

const tokenKey = "vgin:auth-token"
const tokenManagerKey = "vgin:token-manager"

var ()

func TokenManagerMiddleware(tokenTmpl vauth.Token, tokenManager *vauth.TokenManager) gin.HandlerFunc {
	t := reflect.TypeOf(tokenTmpl)
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		lg.Fatal("NewTaskQueue: typeObj should be ptr to struct")
	}

	t = t.Elem()

	return func(c *gin.Context) {
		headerToken := c.GetHeader(AuthHeaderKey)

		token := reflect.New(t).Interface().(vauth.Token)
		token.SetKey(headerToken)

		if err := tokenManager.Read(token); err != nil {
			// no token
			return
		}

		c.Set(tokenManagerKey, tokenManager)
		SetToken(c, token)
	}
}

func TokenRequiredMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := GetToken(c)
		if token == nil {
			vgin.AbortWithError(c, http.StatusUnauthorized, "token required")
			return
		}
		c.Next()
	}
}

func SaveToken(c *gin.Context, token vauth.Token) error {
	m, exists := c.Get(tokenManagerKey)
	if !exists {
		return errors.New("token manager not init")
	}

	tm := m.(*vauth.TokenManager)
	return tm.Save(token)
}

func GetToken(c *gin.Context) vauth.Token {
	val, exists := c.Get(tokenKey)
	if !exists {
		return nil
	}

	return val.(vauth.Token)
}

func SetToken(c *gin.Context, token vauth.Token) {
	c.Set(tokenKey, token)
}
