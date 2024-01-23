package idgenerate

import (
	"time"
)

const (
	maxNext = 0b11111_11111111_111
	shardId = 4
)

var (
	Offset    = str2time("2000-01-01 08:00:00")
	offset    int64
	lastEpoch int64
)

func GetId() int64 {
	i := time.Now().UnixMilli() / 1000
	return nextId(i)
}
func nextId(epochSecond int64) (id int64) {
	if epochSecond < lastEpoch {
		epochSecond = lastEpoch
	}
	if lastEpoch != epochSecond {
		lastEpoch = epochSecond
		offset = 0
	}
	offset++
	next := offset & maxNext
	if next == 0 {
		return nextId(epochSecond + 1)
	}
	return generateId(epochSecond, next, shardId)
}

func generateId(epochSecond, next, shardId int64) int64 {
	return ((epochSecond - Offset) << 21) | (next << 5) | shardId
}

// 日期转化为时间戳
func str2time(datetime string) int64 {
	timeLayout := time.DateTime          //转化所需模板
	loc, _ := time.LoadLocation("Local") //获取时区
	tmp, _ := time.ParseInLocation(timeLayout, datetime, loc)
	return tmp.Unix() //转化¬为时间戳 类型是int64
}
