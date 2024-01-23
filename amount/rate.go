package amount

import "github.com/shopspring/decimal"

// Rate 声明一个金额类型
type Rate string

const Default = 1
const Hundred = 100

// ToString 带 % 的字符串
// base 是否乘基数100
func (r Rate) ToString(base int, defaultStr ...string) string {
	rate := string(r)
	if len(rate) == 0 {
		return defaultStr[0]
	}
	if base == Default {
		return rate + "%"
	}
	float, err := decimal.NewFromString(rate)
	if err != nil {
		return defaultStr[0]
	}
	tmp := decimal.NewFromFloat(100)
	return float.Mul(tmp).String() + "%"
}
