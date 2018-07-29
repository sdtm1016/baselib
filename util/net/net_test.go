package net_util

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
	"baselib/constants"
)

func TestGetLocalIp(t *testing.T) {
	ips, err := GetLocalIp("以太网")
	assert.Nil(t, err)
	fmt.Println("net ip:", ips)

	fmt.Printf("ZERO:%d,ONE:%d,TWO:%d",constants.ZERO,constants.ONE,constants.TWO)
}
