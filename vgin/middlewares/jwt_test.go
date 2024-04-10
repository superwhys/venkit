package middlewares

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/superwhys/venkit/lg"

	"github.com/superwhys/venkit/vhttp"
)

func generateJwtToken() string {
	token, err := GenerateJWTAuth("test-key", &UserInfoClaims{
		User: "yong",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "yong-project",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * 70)),
		},
	})
	lg.PanicError(err)
	return token
}

func TestJwtAuth(t *testing.T) {
	client := vhttp.Default()
	resp := client.Post(context.TODO(), "http://localhost:8081/test_jwt", nil, vhttp.NewHeader().Add(AuthHeaderKey, fmt.Sprintf("Bearer %v", generateJwtToken())))

	respStr, err := resp.BodyString()
	if err != nil {
		t.Error(err)
		return
	}
	if strings.Contains(respStr, "Authorization failure") {
		t.Error("expect success but get failed")
	}
	fmt.Println(respStr)

	resp = client.Post(context.TODO(), "http://localhost:8081/test_jwt", nil, vhttp.NewHeader().Add(AuthHeaderKey, fmt.Sprintf("Bearer %v", "1234")))

	respStr, err = resp.BodyString()
	if err != nil {
		t.Error(err)
		return
	}
	if !strings.Contains(respStr, "Authorization failure") {
		t.Error("expect failed but get success")
	}
	fmt.Println(respStr)

}
