package benchmarks

import (
	"github.com/shirou/gopsutil/cpu"
	"fmt"
	"github.com/shirou/gopsutil/mem"
	"strings"
	"github.com/shirou/gopsutil/process"
	"time"
)

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

func CpuCoreCount() (int) {
	count, _ := cpu.Counts(true)
	return count
}

func CpuUsageAsFloat(proc *process.Process) float64 {
	usage, _ := proc.Percent(time.Duration(200) * time.Millisecond)
	return usage
}

func CpuUsageAsString(proc *process.Process) string {
	return FloatToString(CpuUsageAsFloat(proc))
}

func MemorySize() string {
	memory, _ := mem.VirtualMemory()
	return ByteCountDecimal(memory.Total)
}

func MemoryUsageAsFloat(proc *process.Process) float64 {
	usage, _ := proc.MemoryPercent()
	return float64(usage)
}

func MemoryUsageAsString(proc *process.Process) string {
	return FloatToString(MemoryUsageAsFloat(proc))
}

func FloatToString(value float64) string {
	return strings.Replace(fmt.Sprintf("%.2f%%", value), ".", ",", -1)
}
