package hardware

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/load" // CPU负载
	"github.com/shirou/gopsutil/v3/mem"  // 内存占用
	"runtime"
	"time"
)

type Data struct {
	MemTotal     uint64  // 内存大小
	MemAvailable uint64  // 闲置可用内存
	MemFree      uint64  // 未使用内存
	MemUsed      uint64  // 已使用内存
	UsedPercent  float64 // 已使用百分比
	CpuNum       int     // 逻辑CPU数量
	CpuPercent   float64 // CPU使用率
	CpuLoad1     float64 // 每1分钟平均负载
	CpuLoad5     float64 // 每5分钟平均负载
	CpuLoad15    float64 // 每15分钟平均负载
}

func Hardware() (data Data) {
	// 内存信息
	memVirtual, _ := mem.VirtualMemory()
	data.MemTotal = memVirtual.Total / 1024 / 1024
	data.MemAvailable = memVirtual.Available / 1024 / 1024
	data.MemFree = memVirtual.Free / 1024 / 1024
	data.MemUsed = memVirtual.Used / 1024 / 1024
	data.UsedPercent = memVirtual.UsedPercent / 1024 / 1024
	// CPU信息
	data.CpuNum = runtime.NumCPU()
	cpuPercent, _ := cpu.Percent(time.Second, false)
	data.CpuPercent = cpuPercent[0]
	// CPU负载（不耗时）
	cpuLoad, _ := load.Avg()
	// {"load1":3.62109375,"load5":2.93408203125,"load15":2.58251953125}
	// load表示每1分钟、5分钟、15分钟的平均队列（平均负载）,值为进程或线程数
	// 具体示意请参考load average：https://blog.csdn.net/bd_zengxinxin/article/details/51781630
	data.CpuLoad1 = cpuLoad.Load1
	data.CpuLoad5 = cpuLoad.Load5
	data.CpuLoad15 = cpuLoad.Load15
	return data
}
