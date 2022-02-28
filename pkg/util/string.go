package util

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

// MaskMobile 隐藏手机号中间四位数字
func MaskMobile(str string) (result string) {
	if str == "" {
		return ""
	}
	err := MatchMobile(str)
	if err == nil {
		result = Substr2(str, 0, 3) + "****" + Substr2(str, 7, 11)
	} else {
		result = str
		//nameRune := []rune(str)
		//lens := len(nameRune)
		//if lens <= 1 {
		//	result = "***"
		//} else if lens == 2 {
		//	result = string(nameRune[:1]) + "*"
		//} else if lens == 3 {
		//	result = string(nameRune[:1]) + "*" + string(nameRune[2:3])
		//} else if lens == 4 {
		//	result = string(nameRune[:1]) + "**" + string(nameRune[lens-1:lens])
		//} else if lens > 4 {
		//	result = string(nameRune[:2]) + "***" + string(nameRune[lens-2:lens])
		//}
	}
	return
}

// MatchMobile 正则校验手机号
func MatchMobile(mobile string) error {
	reg := `^1[0-9]\d{9}$`
	rgx := regexp.MustCompile(reg)
	mobileMatch := rgx.MatchString(mobile)
	if mobileMatch {
		return nil
	}
	return errors.New("手机号码有误")
}
func Substr2(str string, start int, end int) string {
	rs := []rune(str)
	return string(rs[start:end])
}

func MD5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}

// StringToInt string转int
func StringToInt(_str string) int64 {
	_int, err := strconv.ParseInt(_str, 10, 64) // string转int
	if err != nil {                             // 报错则默认返回0
		_int = 0
		//fmt.Println("格式转换错误，默认为0。")
		//fmt.Println(err)
	}
	return _int
}

// IntToString int转string
func IntToString(_int int64) string {
	_str := strconv.FormatInt(_int, 10)
	return _str
}

// InterfaceToString 获取变量的字符串值
// 浮点型 3.0将会转换成字符串3, "3"
// 非数值或字符类型的变量将会被转换成JSON格式字符串
func InterfaceToString(value interface{}) string {
	var key string
	if value == nil {
		return key
	}

	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}

	return key
}

// Is c an ASCII lower-case letter?
func isASCIILower(c byte) bool {
	return 'a' <= c && c <= 'z'
}

// Is c an ASCII digit?
func isASCIIDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

// CamelCase returns the CamelCased name.
// If there is an interior underscore followed by a lower case letter,
// drop the underscore and convert the letter to upper case.
// There is a remote possibility of this rewrite causing a name collision,
// but it's so remote we're prepared to pretend it's nonexistent - since the
// C++ generator lowercases names, it's extremely unlikely to have two fields
// with different capitalizations.
// In short, _my_field_name_2 becomes XMyFieldName_2.
func CamelCase(s string) string {
	if s == "" {
		return ""
	}
	t := make([]byte, 0, 32)
	i := 0
	if s[0] == '_' {
		// Need a capital letter; drop the '_'.
		t = append(t, 'X')
		i++
	}
	// Invariant: if the next letter is lower case, it must be converted
	// to upper case.
	// That is, we process a word at a time, where words are marked by _ or
	// upper case letter. Digits are treated as words.
	for ; i < len(s); i++ {
		c := s[i]
		if c == '_' && i+1 < len(s) && isASCIILower(s[i+1]) {
			continue // Skip the underscore in s.
		}
		if isASCIIDigit(c) {
			t = append(t, c)
			continue
		}
		// Assume we have a letter now - if not, it's a bogus identifier.
		// The next word is a sequence of characters that must start upper case.
		if isASCIILower(c) {
			c ^= ' ' // Make it a capital letter.
		}
		t = append(t, c) // Guaranteed not lower case.
		// Accept lower case sequence that follows.
		for i+1 < len(s) && isASCIILower(s[i+1]) {
			i++
			t = append(t, s[i])
		}
	}
	return string(t)
}
