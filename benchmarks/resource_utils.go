package benchmarks

import (
	"github.com/shirou/gopsutil/cpu"
	"fmt"
	"github.com/shirou/gopsutil/mem"
	"strings"
	"github.com/shirou/gopsutil/process"
)

func currentProcess() *process.Process {
	processes, _ := process.Processes()
	return processes[0]
}

func CpuCoreCount() (int) {
	count, _ := cpu.Counts(true)
	return count
}

func CpuUsageAsFloat() float64 {
	usage, _ := currentProcess().CPUPercent()
	return usage
}

func CpuUsageAsString() string {
	return FloatToString(CpuUsageAsFloat())
}

func ByteCountDecimal(b uint64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}

func MemorySize() string {
	memory, _ := mem.VirtualMemory()
	return ByteCountDecimal(memory.Total)
}

func MemoryUsageAsFloat() float64 {
	usage, _ := currentProcess().MemoryPercent()
	return float64(usage)
}

func MemoryUsageAsString() string {
	return FloatToString(MemoryUsageAsFloat())
}

func FloatToString(value float64) string {
	return strings.Replace(fmt.Sprintf("%.2f%%", value), ".", ",", -1)
}
