package sign_util

import (
	"fmt"
	"github.com/golang-module/dongle"
	"net/url"
	"sort"
	"strconv"
	"time"
)

// MakeSign 接口请求加签
func MakeSign(data url.Values, signKey string) url.Values {
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	data.Set("timestamp", timestamp)
	str1 := kSoftValues(data)
	sign := dongle.Encrypt.FromString(str1 + signKey).ByMd5().String()
	data.Set("sign", sign)
	return data
}

// CheckSign 验证签名
func CheckSign(data url.Values, signKey string) bool {
	oldSign := data.Get("sign")
	data.Del("sign")
	s := kSoftValues(data)
	newSign := dongle.Encrypt.FromString(s + signKey).ByMd5().String()
	if oldSign == newSign {
		return true
	}
	return false
}

// 签名参数排序
func kSoftValues(values url.Values) string {
	var slice []string
	for key := range values {
		slice = append(slice, key)
	}
	sort.Strings(slice)
	var s string
	index := 0
	for _, v := range slice {
		sprintf := fmt.Sprintf("%s=%s", v, values.Get(v))
		if index == 0 {
			s = sprintf
		} else {
			s = s + "&" + sprintf
		}
		index++
	}
	return s
}
