package utils

import (
	"strconv"
	"time"
)

// 统计字符串长度
func StrLen(str string) int {
	return len([]rune(str))
}

// 时间戳毫秒 -> format time(Local)。
// if rawTime is "" , return current time。
func GetFormatedTimeFromUnixMilliStr(rawTime string, format string) (string, error) {
	if rawTime == "" || rawTime == "0" {
		return time.Now().Format(format), nil
	}
	// 时间戳秒 -> format time
	latestTime, err := strconv.ParseInt(rawTime, 10, 64)
	if err != nil {
		return "", err
	}
	foramtedTime := time.UnixMilli(int64(latestTime)).In(time.Local).Format(format)
	return foramtedTime, nil
}

// 时间戳毫秒 -> format time(Local)。
// if rawTime is "" , return current time。
func GetFormatedTimeFromUnixMilli(rawTime int64, format string) string {
	return time.UnixMilli(rawTime).In(time.Local).Format(format)
}
