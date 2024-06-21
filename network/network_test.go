package network

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMacAddr(t *testing.T) {
	mac, err := GetMacAddr("en8")
	assert.Nil(t, err)
	fmt.Println(mac)
}

func TestGetLocalIp(t *testing.T) {
	ip, err := GetLocalIP("en8")
	assert.Nil(t, err)
	fmt.Println(ip)
}

func TestGetSubnetMask(t *testing.T) {
	sm, err := GetSubnetMask("en8")
	assert.Nil(t, err)
	fmt.Println(sm)
}
