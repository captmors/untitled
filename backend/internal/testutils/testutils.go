package testutils

import (
	"fmt"
	"sync/atomic"
)

const basePort = 1488

var uniqueNum int32

func GetUniquePort() string {
	uniquePort := atomic.AddInt32(&uniqueNum, 1) + int32(basePort)
	return fmt.Sprintf("%d", uniquePort)
}
