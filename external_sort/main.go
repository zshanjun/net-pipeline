package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"zshanjun/test/pipeline"
)

func main() {
	p := createNetworkPipeline("large.in", 800000000, 4)
	writeToFile(p, "large.out")
	printFile("large.out")
}

func printFile(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	p := pipeline.ReaderSource(bufio.NewReader(file), -1)
	count := 0
	for v := range p {
		fmt.Println(v)
		count++
		if count >= 100 {
			break
		}
	}
}

func writeToFile(in <-chan int, fileName string) {
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	pipeline.WriterSink(writer, in)
}

func createPipeline(fileName string, fileSize, chunkCount int) <-chan int {
	chunkSize := fileSize / chunkCount
	sortResults := []<-chan int{}
	pipeline.Init()
	for i := 0; i < chunkCount; i++ {
		file, err := os.Open(fileName)
		if err != nil {
			panic(err)
		}
		file.Seek(int64(i*chunkSize), 0)
		source := pipeline.ReaderSource(bufio.NewReader(file), chunkSize)
		sortResults = append(sortResults, pipeline.InMemorySort(source))
	}

	return pipeline.MergeN(sortResults...)
}

func createNetworkPipeline(fileName string, fileSize, chunkCount int) <-chan int {
	chunkSize := fileSize / chunkCount
	sortAddr := []string{}
	pipeline.Init()
	for i := 0; i < chunkCount; i++ {
		file, err := os.Open(fileName)
		if err != nil {
			panic(err)
		}
		file.Seek(int64(i*chunkSize), 0)
		source := pipeline.ReaderSource(bufio.NewReader(file), chunkSize)
		addr := ":" + strconv.Itoa(7000+i)
		pipeline.NetworkSink(addr, pipeline.InMemorySort(source))
		sortAddr = append(sortAddr, addr)
	}

	sortResults := []<-chan int{}
	for _, addr := range sortAddr {
		sortResults = append(sortResults, pipeline.NetworkSource(addr))
	}

	return pipeline.MergeN(sortResults...)
}
