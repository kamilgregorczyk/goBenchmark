package main

import (
	"fmt"
	"io/ioutil"
	"encoding/json"
	"regexp"
	"strconv"
	"goBenchmark/schema"
	"goBenchmark/benchmarks"
)

func LoadJsonSchemaFromFileAsStruct() (schema.JsonSchema) {
	byteArray, e := ioutil.ReadFile("test.json")
	if e != nil {
		panic(e)
	}
	var jsonSchema schema.JsonSchema
	json.Unmarshal(byteArray, &jsonSchema)
	return jsonSchema

}

func LoadJsonSchemaFromFileAsString() (string) {
	byteArray, e := ioutil.ReadFile("test.json")
	if e != nil {
		panic(e)
	}
	return string(byteArray)
}

func PrepareCSV() {
	e := ioutil.WriteFile("results.csv", []byte("Benchmark name;runtime;cpu usage;memory usage\n"), 0777)
	if e != nil {
		panic(e)
	}
}

func UnorderedListOfNumbers(value string) ([]int) {
	r, _ := regexp.Compile("[+-]?[0-9]+")
	values := r.FindAllString(value, -1)
	var newValues []int
	for _, i := range values {
		int_value, e := strconv.Atoi(i)
		if e != nil {
			panic(e)
		}
		newValues = append(newValues, int_value)
	}

	return newValues[:(len(newValues) / 11)]
}

func main() {
	fmt.Println("************************")
	fmt.Println("***   Go benchmark   ***")
	fmt.Println("************************")
	fmt.Println("Available CPU cores: " + fmt.Sprintf("%v", benchmarks.CpuCoreCount()) + " Current CPU usage: " + benchmarks.CpuUsageAsString())
	fmt.Println("Available Memory: " + benchmarks.MemorySize() + " Current Memory usage: " + benchmarks.MemoryUsageAsString())
	testDataAsJsonChan := make(chan schema.JsonSchema, 1)
	testDataAsStringChan := make(chan string, 1)
	unorderedListOfNumbersChan := make(chan []int, 1)
	go func() {
		PrepareCSV()
	}()
	go func() {
		testDataAsJsonChan <- LoadJsonSchemaFromFileAsStruct()
		close(testDataAsJsonChan)
	}()
	go func() {
		dataAsString := LoadJsonSchemaFromFileAsString()
		testDataAsStringChan <- dataAsString
		go func() {
			unorderedListOfNumbersChan <- UnorderedListOfNumbers(dataAsString)
			close(unorderedListOfNumbersChan)
		}()
		close(testDataAsStringChan)
	}()
	testDataAsJson := <-testDataAsJsonChan
	testDataAsString := <-testDataAsStringChan
	unorderedListOfNumbers := <-unorderedListOfNumbersChan

	benchmarks.BTreeBenchmark(benchmarks.NewTimer("B tree", testDataAsString, testDataAsJson, unorderedListOfNumbers))
	benchmarks.BuiltInSortBenchmark(benchmarks.NewTimer("Built-in sort", testDataAsString, testDataAsJson, unorderedListOfNumbers))
	benchmarks.MergeSortBenchmark(benchmarks.NewTimer("Merge sort", testDataAsString, testDataAsJson, unorderedListOfNumbers))
	benchmarks.RegexpBenchmark(benchmarks.NewTimer("Regexp for digits", testDataAsString, testDataAsJson, unorderedListOfNumbers))
	benchmarks.JsonImportBenchmark(benchmarks.NewTimer("Importing big json file", testDataAsString, testDataAsJson, unorderedListOfNumbers))
	benchmarks.AggregateColumnBenchmark(benchmarks.NewTimer("Aggregating column from dict and counting median", testDataAsString, testDataAsJson, unorderedListOfNumbers))

	fmt.Println("Done!")

}
