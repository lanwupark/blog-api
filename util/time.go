package util

import (
	"fmt"
	"time"
)

// GetIntervalString 获取两个时间间隔
func GetIntervalString(from time.Time, to time.Time) string {
	if from.After(to) {
		from, to = to, from //交换
	}
	if interval := to.Year() - from.Year(); interval > 0 {
		return fmt.Sprintf("%d年前", interval)
	}
	if interval := to.Month() - from.Month(); interval > 0 {
		return fmt.Sprintf("%d月前", interval)
	}
	if interval := to.Day() - from.Day(); interval > 0 {
		return fmt.Sprintf("%d天前", interval)
	}
	if interval := to.Hour() - from.Hour(); interval > 0 {
		return fmt.Sprintf("%d小时前", interval)
	}
	if interval := to.Minute() - from.Minute(); interval > 0 {
		return fmt.Sprintf("%d分钟前", interval)
	}
	return fmt.Sprintf("%d秒前", to.Second()-from.Second())
}
