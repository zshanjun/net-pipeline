package main

import (
	"bufio"
	"fmt"
	"os"
	"zshanjun/net-pipeline/pipeline"
)

func main() {
	var filename = "large.in"
	var count = 100000000
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	p := pipeline.RandomSource(count)
	writer := bufio.NewWriter(file)
	pipeline.WriterSink(writer, p)
	writer.Flush()

	file, err = os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	p = pipeline.ReaderSource(bufio.NewReader(file), -1)
	counter := 0
	for v := range p {
		fmt.Println(v)
		counter++
		if counter > 100 {
			break
		}
	}
}

func mergeDemo() {
	p := pipeline.Merge(
		pipeline.InMemorySort(pipeline.ArraySource(3, 1, 2, 5, 9, 4)),
		pipeline.InMemorySort(pipeline.ArraySource(10, 4, 92, 59, 2, 0)),
	)
	for v := range p {
		fmt.Println(v)
	}
}
