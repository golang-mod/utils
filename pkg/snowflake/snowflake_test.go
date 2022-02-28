package snowflake

import (
	"fmt"
	"testing"
)

func TestId(t *testing.T) {
	// 生成节点实例
	node, _ := New(1)
	got := node.GetId()
	fmt.Println(got)
	if got == 0 {
		t.Errorf("expected") // 测试失败输出错误提示
	}
}

// 测试函数的性能
//  go test -bench=Id
//  go test -bench=Id -benchmem
func BenchmarkId(b *testing.B) {
	node, _ := New(1)
	// b.N不是固定的数
	for i := 0; i < b.N; i++ {
		node.GetId()
	}
}
