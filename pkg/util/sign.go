package util

import (
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"time"
)

// MakeSign 接口请求加签
func MakeSign(data url.Values, signKey string) url.Values {
	timestamp := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	data.Set("timestamp", timestamp)
	str := kSoftValues(data)
	sign := MD5(str + signKey)
	data.Set("sign", sign)
	return data
}

// CheckSign 验证签名
func CheckSign(data url.Values, signKey string) bool {
	oldSign := data.Get("sign")
	data.Del("sign")
	str := kSoftValues(data)
	newSign := MD5(str + signKey)
	if oldSign == newSign {
		return true
	}
	return false
}

// 签名参数排序
func kSoftValues(values url.Values) string {
	var sslice []string
	for key, _ := range values {
		sslice = append(sslice, key)
	}
	sort.Strings(sslice)
	var str string
	index := 0
	for _, v := range sslice {
		sprintf := fmt.Sprintf("%s=%s", v, values.Get(v))
		if index == 0 {
			str = sprintf
		} else {
			str = str + "&" + sprintf
		}
		index++
	}
	return str
}
