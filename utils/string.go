package utils

import (
	"strconv"
	"time"
)

// 统计字符串长度
func StrLen(str string) int {
	return len([]rune(str))
}

// getFormatedTime。
// 时间戳秒 -> format time。
// if rawTime is "" , return current time。
func GetFormatedTimeFromUnix(rawTime string, format string) (string, error) {
	if rawTime == "" {
		return time.Now().Format(format), nil
	}
	// 时间戳秒 -> format time
	latestTime, err := strconv.ParseInt(rawTime, 10, 64)
	if err != nil {
		return "", err
	}
	foramtedTime := time.Unix(int64(latestTime), 0).Format(format)
	return foramtedTime, nil
}
