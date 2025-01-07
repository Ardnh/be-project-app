package helper

import (
	"fmt"
	"time"
)

func StartTime() time.Time {
	return time.Now()
}

func ElapsedTime(t time.Time) string {
	elapsed := time.Since(t).Milliseconds()
	elapsedString := fmt.Sprintf("%d ms", elapsed)
	return elapsedString
}
