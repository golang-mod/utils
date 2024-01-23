package dbgenerate

import "testing"

// 生成数据表结构体
func TestGenerate(t *testing.T) {
	Generate(nil, Config{
		Path:        "./models",
		Tables:      nil,
		TablePrefix: "",
	})
}
