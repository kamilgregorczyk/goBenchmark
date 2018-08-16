package benchmarks

import (
	"sort"
	"regexp"
	"io/ioutil"
	"goBenchmark/schema"
	"encoding/json"
	"strconv"
	"github.com/montanaflynn/stats"
)

func BTreeBenchmark(timer Timer) {
	timer.startBenchmark()
	defer timer.stopBenchmark()
	tree := &BinaryTree{}
	i := 0
	for {
		if (i >= 10000) {
			break
		}
		tree.insert(i)
		i++
	}

}
func MergeSortBenchmark(timer Timer) {
	timer.startBenchmark()
	defer timer.stopBenchmark()
	Mergesort(timer.unorderedListOfNumbers)

}

func BuiltInSortBenchmark(timer Timer) {
	timer.startBenchmark()
	defer timer.stopBenchmark()
	sort.Ints(timer.unorderedListOfNumbers)
}

func RegexpBenchmark(timer Timer) {
	timer.startBenchmark()
	defer timer.stopBenchmark()
	r, _ := regexp.Compile("[+-]?[0-9]+")
	r.FindAllString(timer.dataAsString, -1)
}

func JsonImportBenchmark(timer Timer) {
	timer.startBenchmark()
	defer timer.stopBenchmark()
	byteArray, e := ioutil.ReadFile("test.json")
	if e != nil {
		panic(e)
	}
	var jsonSchema schema.JsonSchema
	json.Unmarshal(byteArray, &jsonSchema)
}

func AggregateColumnBenchmark(timer Timer) {
	timer.startBenchmark()
	defer timer.stopBenchmark()
	iterations := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	for range iterations {
		var from_values []int64
		var to_values []int64
		var result []float64
		for _, feature := range timer.dataAsJson.Features {
			if feature.Properties.FROMST != "" {
				value, err := strconv.ParseInt(feature.Properties.FROMST, 10, 64)
				if err != nil {
					panic(err)
				}
				from_values = append(from_values, value)
			}
		}
		for _, feature := range timer.dataAsJson.Features {
			if feature.Properties.TOST != "" {
				value, err := strconv.ParseInt(feature.Properties.TOST, 10, 64)
				if err != nil {
					panic(err)
				}
				to_values = append(to_values, value)
			}
		}
		for index, _ := range from_values {
			result = append(result, float64(from_values[index]*to_values[index]))
		}
		stats.Median(result)
	}

}
