package vauth

import (
	"errors"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/superwhys/venkit/lg"
	"github.com/superwhys/venkit/vgin/v2"
)

const (
	tokenKey        = "vgin:auth-token"
	tokenContextKey = "vgin:token-manager"
	tokenTemplate   = "vgin:token-template"
	AuthHeaderKey   = "Authorization"
)

type tokenContext struct {
	tokenManager *TokenManager
	tokenTmpl    reflect.Type
}

func (t *tokenContext) TokenManager() *TokenManager {
	return t.tokenManager
}

func TokenManagerMiddleware(tokenTmpl Token, tokenManager *TokenManager) gin.HandlerFunc {
	t := reflect.TypeOf(tokenTmpl)
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		lg.Fatal("TokenManagerMiddleware: token template should be ptr to struct")
	}

	t = t.Elem()

	return func(c *gin.Context) {
		tokenCtx := &tokenContext{
			tokenManager: tokenManager,
			tokenTmpl:    t,
		}
		c.Set(tokenContextKey, tokenCtx)
	}
}

func CurrentTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader(AuthHeaderKey)
		if tokenStr == "" {
			return
		}
		tokenCtx := getTokenContext(c)
		if tokenCtx == nil {
			return
		}
		newToken := genNewToken(c)
		if newToken == nil {
			return
		}

		if err := tokenCtx.TokenManager().Read(tokenStr, newToken); err != nil {
			lg.Errorf("token manager read token: %v error: %v", tokenStr, err)
			return
		}

		SetToken(c, newToken)
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

func SaveToken(c *gin.Context, token Token) error {
	m, exists := c.Get(tokenContextKey)
	if !exists {
		return errors.New("token manager not init")
	}

	ctx := m.(*tokenContext)
	return ctx.tokenManager.Save(token)
}

func GetToken(c *gin.Context) Token {
	val, exists := c.Get(tokenKey)
	if !exists {
		return nil
	}

	return val.(Token)
}

func SetToken(c *gin.Context, token Token) {
	c.Set(tokenKey, token)
}

func getTokenContext(c *gin.Context) *tokenContext {
	tc, exists := c.Get(tokenContextKey)
	if !exists {
		return nil
	}

	return tc.(*tokenContext)
}

func genNewToken(c *gin.Context) Token {
	tt, exists := c.Get(tokenTemplate)
	if !exists {
		return nil
	}
	tokenTemplate := tt.(reflect.Type)
	return reflect.New(tokenTemplate).Interface().(Token)
}
