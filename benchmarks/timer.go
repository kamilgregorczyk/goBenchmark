package benchmarks

import (
	"time"
	"goBenchmark/schema"
	"fmt"
	"strings"
	"os"
	"github.com/shirou/gopsutil/process"
)

type Metric struct {
	cpuUsagePoints    []float64
	memoryUsagePoints []float64
	stopEvent         chan bool
	stoppedEvent      chan bool
}
type Timer struct {
	name                   string
	dataAsString           string
	dataAsJson             schema.JsonSchema
	unorderedListOfNumbers []int
	start                  time.Time
	runtime                time.Duration
	runtimeFormatted       string
	cpuUsage               string
	memoryUsage            string
	metrics                Metric
	proc                   *process.Process
}

func NewTimer(name string, dataAsString string, dataAsJson schema.JsonSchema, unorderedListOfNumbers []int, proc *process.Process) Timer {
	timer := Timer{}
	timer.name = name
	timer.dataAsString = dataAsString
	timer.dataAsJson = dataAsJson
	timer.unorderedListOfNumbers = unorderedListOfNumbers
	timer.metrics = Metric{stopEvent: make(chan bool), stoppedEvent: make(chan bool)}
	timer.proc = proc
	return timer
}

func (timer *Timer) startBenchmark() {
	fmt.Println()
	fmt.Println("Starting '" + timer.name + "' benchmark...")
	timer.startTimer()
	timer.startRecordingResources()
}
func (timer *Timer) stopBenchmark() {
	timer.stopTimer()
	timer.stopRecordingResources()
	timer.runtimeFormatted = strings.Replace(fmt.Sprintf("%.2f", timer.runtime.Seconds()), ".", ",", -1)
	fmt.Println("Benchmark took " + timer.runtimeFormatted + " seconds, used " + timer.cpuUsage + " CPU and " + timer.memoryUsage + " RAM")
	timer.saveResults()
}

func (timer *Timer) startRecordingResources() {
	timer.recordUsage()
	go func() {
		for {
			select {
			case <-timer.metrics.stopEvent:
				timer.metrics.stoppedEvent <- true
				return
			default:
				timer.recordUsage()
			}
		}
	}()
}

func (timer *Timer) recordUsage() {
	timer.metrics.cpuUsagePoints = append(timer.metrics.cpuUsagePoints, CpuUsageAsFloat(timer.proc))
	timer.metrics.memoryUsagePoints = append(timer.metrics.memoryUsagePoints, MemoryUsageAsFloat(timer.proc))
}
func (timer *Timer) stopRecordingResources() {
	timer.metrics.stopEvent <- true
	<-timer.metrics.stoppedEvent
	timer.cpuUsage = timer.calculateCpuUsage()
	timer.memoryUsage = timer.calculateMemoryUsage()
}

func (timer *Timer) calculateCpuUsage() string {
	var total float64 = 0
	for _, value := range timer.metrics.cpuUsagePoints {
		total += value
	}
	return FloatToString(total / float64(len(timer.metrics.cpuUsagePoints)))
}

func (timer *Timer) calculateMemoryUsage() string {
	var total float64 = 0
	for _, value := range timer.metrics.memoryUsagePoints {
		total += value
	}
	return FloatToString(total / float64(len(timer.metrics.memoryUsagePoints)))
}

func (timer *Timer) startTimer() {
	timer.start = time.Now()
}
func (timer *Timer) stopTimer() {
	end := time.Now()
	timer.runtime = end.Sub(timer.start)
}
func (timer *Timer) saveResults() {
	line := fmt.Sprintf("%s;%s;%s;%s\n", timer.name, timer.runtimeFormatted, timer.cpuUsage, timer.memoryUsage)
	file, err := os.OpenFile("results.csv", os.O_APPEND|os.O_WRONLY, 0600)
	defer file.Close()
	if err != nil {
		panic(err)
	}
	fmt.Fprint(file, line)

}
