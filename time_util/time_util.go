package time_util

import "fmt"

// ResolveTime 秒转换为时分秒
func ResolveTime(seconds int64) (duration string) {
	var day = seconds / (24 * 3600)
	hour := (seconds - day*3600*24) / 3600
	minute := (seconds - day*24*3600 - hour*3600) / 60
	second := seconds - day*24*3600 - hour*3600 - minute*60
	if hour > 0 {
		duration = fmt.Sprintf("%d%s", hour, "小时")
	}
	if minute > 0 {
		duration = fmt.Sprintf("%s%d%s", duration, minute, "分钟")
	}
	if second > 0 {
		duration = fmt.Sprintf("%s%d%s", duration, second, "秒")
	}
	return duration
}
