package vgin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type headerData struct {
	Token  string `vheader:"token"`
	UserId int    `vheader:"user_id"`
}

func TestMapHeader(t *testing.T) {
	h := map[string][]string{
		"token":   {"test-token"},
		"user_id": {"1111"},
	}

	hd := new(headerData)
	err := mapFormByTag(hd, h, "vheader")
	assert.Nil(t, err)
	assert.Equal(t, hd, &headerData{Token: "test-token", UserId: 1111})
}
