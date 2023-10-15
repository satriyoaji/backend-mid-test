package util

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateAccountNumber(max, min int) string {
	rand.Seed(time.Now().UnixNano())
	randomNum := rand.Intn(max-min+1) + min

	return fmt.Sprintf("%d", randomNum)
}
