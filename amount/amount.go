package amount

import "github.com/shopspring/decimal"

// Amount 声明一个金额类型
type Amount float64

// WanToFen 万转分
func (a Amount) WanToFen() int64 {
	float := decimal.NewFromFloat(float64(a))
	tmp := decimal.NewFromFloat(1000000)
	return float.Mul(tmp).IntPart()
}
func (a Amount) WanToFenPointer() *int64 {
	var t = a.WanToFen()
	return &t
}

// YuanToFen 元转分
func (a Amount) YuanToFen() int64 {
	float := decimal.NewFromFloat(float64(a))
	tmp := decimal.NewFromFloat(100)
	return float.Mul(tmp).IntPart()
}
func (a Amount) YuanToFenPointer() *int64 {
	var t = a.YuanToFen()
	return &t
}

// FenToYuan 分转元
func (a Amount) FenToYuan() float64 {
	float := decimal.NewFromFloat(float64(a))
	tmp := decimal.NewFromFloat(100)
	f, _ := float.Div(tmp).Float64()
	return f
}
func (a Amount) FenToYuanPointer() *float64 {
	var t = a.FenToYuan()
	return &t
}

// FenToWan 分转万
func (a Amount) FenToWan() float64 {
	float := decimal.NewFromFloat(float64(a))
	tmp := decimal.NewFromFloat(1000000)
	f, _ := float.Div(tmp).Float64()
	return f
}
func (a Amount) FenToWanPointer() *float64 {
	var t = a.FenToWan()
	return &t
}
