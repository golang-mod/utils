package time_util

import (
	"math"
	"strconv"
	"strings"
	"time"
)

// GetTimeDate 获取日期时间戳，s
// Y年m月d号 H:i:s.MS.NS 星期W
func GetTimeDate(_format string) (date string) {
	if len(_format) == 0 {
		_format = "YmdHisMS"
	}
	date = _format

	// 时区
	//timeZone, _ := time.LoadLocation(ServerInfo["timezone"])
	timeZone := time.FixedZone("CST", 8*3600) // 东八区

	timer := time.Now().In(timeZone)

	var year = int64(timer.Year())
	var month = int64(timer.Month())
	var day = int64(timer.Day())
	var hour = int64(timer.Hour())
	var minute = int64(timer.Minute())
	var second = int64(timer.Second())
	var week = int64(timer.Weekday())
	var ms = timer.UnixNano() / 1e6
	var ns = timer.UnixNano() / 1e9
	msTmp := intToTime(int64(math.Floor(float64(ms / 1000))))
	nsTmp := intToTime(int64(math.Floor(float64(ns / 1000000))))

	var _year string
	var _month string
	var _day string
	var _hour string
	var _minute string
	var _second string
	var _week string // 英文星期
	var _Week string // 中文星期
	var _ms string   // 毫秒
	var _ns string   // 纳秒

	_year = intToTime(year)
	if month < 10 {
		_month = intToTime(month)
		_month = "0" + _month
	} else {
		_month = intToTime(month)
	}
	if day < 10 {
		_day = intToTime(day)
		_day = "0" + _day
	} else {
		_day = intToTime(day)
	}
	if hour < 10 {
		_hour = intToTime(hour)
		_hour = "0" + _hour
	} else {
		_hour = intToTime(hour)
	}
	if minute < 10 {
		_minute = intToTime(minute)
		_minute = "0" + _minute
	} else {
		_minute = intToTime(minute)
	}
	if second < 10 {
		_second = intToTime(second)
		_second = "0" + _second
	} else {
		_second = intToTime(second)
	}
	_week = intToTime(week)
	WeekZh := [...]string{"日", "一", "二", "三", "四", "五", "六"} // 默认从"日"开始
	_Week = WeekZh[week]
	_ms = strings.Replace(intToTime(ms), msTmp, "", -1)
	_ns = strings.Replace(intToTime(ns), nsTmp, "", -1)

	// 替换关键词
	date = strings.Replace(date, "MS", _ms, -1)
	date = strings.Replace(date, "NS", _ns, -1)
	date = strings.Replace(date, "Y", _year, -1)
	date = strings.Replace(date, "m", _month, -1)
	date = strings.Replace(date, "d", _day, -1)
	date = strings.Replace(date, "H", _hour, -1)
	date = strings.Replace(date, "i", _minute, -1)
	date = strings.Replace(date, "s", _second, -1)
	date = strings.Replace(date, "W", _Week, -1)
	date = strings.Replace(date, "w", _week, -1)

	return
}

func intToTime(num int64) string {
	return strconv.Itoa(int(num))
}
