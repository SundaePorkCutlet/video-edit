package utils

import (
	"fmt"
	"time"
)

func GetTime() string {
	now := time.Now()
	return fmt.Sprintf("%d%02d%02d%02d%02d%02d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
}
