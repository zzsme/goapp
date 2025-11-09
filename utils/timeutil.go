// Placeholder for timeutil.go
package utils

import (
	"time"
)

// UnixToTime 格式化时间戳
func UnixToTime(ts int64) string {
	return time.Unix(ts, 0).Format("2006-01-02 15:04:05")
}

// TimeToUnix 时间字符串转时间戳
func TimeToUnix(t string) int64 {
	tm, _ := time.ParseInLocation("2006-01-02 15:04:05", t, time.Local)
	return tm.Unix()
}

// NowTimeString 当前时间字符串
func NowTimeString() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// FormatTime 格式化任意时间对象
func FormatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// DiffDays 计算两个时间相差天数
func DiffDays(t1, t2 time.Time) int {
	return int(t2.Sub(t1).Hours() / 24)
}
