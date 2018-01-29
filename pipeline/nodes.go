package pipeline

import (
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"sort"
	"time"
)

var startTime time.Time

func Init() {
	startTime = time.Now()
}

//数组来源
func ArraySource(num ...int) <-chan int {
	var out = make(chan int)
	go func() {
		for _, v := range num {
			out <- v
		}
		close(out)
	}()
	return out
}

//内存中排序
func InMemorySort(in <-chan int) <-chan int {
	out := make(chan int, 1024)
	go func() {
		//read into memory
		a := []int{}
		for v := range in {
			a = append(a, v)
		}
		fmt.Println("Read done", time.Now().Sub(startTime))

		//sort
		sort.Ints(a)
		fmt.Println("InMemSort done", time.Now().Sub(startTime))

		//send to output
		for _, v := range a {
			out <- v
		}
		close(out)
	}()

	return out
}

//合并
func Merge(in1, in2 <-chan int) <-chan int {
	out := make(chan int, 1024)
	go func() {
		v1, ok1 := <-in1
		v2, ok2 := <-in2
		for ok1 || ok2 {
			if !ok2 || (ok1 && v1 <= v2) {
				out <- v1
				v1, ok1 = <-in1
			} else {
				out <- v2
				v2, ok2 = <-in2
			}
		}
		close(out)
		fmt.Println("Merge done", time.Now().Sub(startTime))
	}()
	return out
}

//实现了reader接口的来源
func ReaderSource(reader io.Reader, chunkSize int) <-chan int {
	out := make(chan int, 1024)
	go func() {
		//64位系统，int为8个字节
		buf := make([]byte, 8)
		readSize := 0
		for {
			n, err := reader.Read(buf)
			readSize += n
			if n > 0 {
				out <- int(binary.BigEndian.Uint64(buf))
			}
			// -1 代表不分割
			if err != nil || (chunkSize != -1 && readSize >= chunkSize) {
				break
			}
		}
		close(out)
	}()
	return out
}

//写入
func WriterSink(writer io.Writer, in <-chan int) {
	for v := range in {
		//64位系统，int为8个字节
		buf := make([]byte, 8)
		binary.BigEndian.PutUint64(buf, uint64(v))
		writer.Write(buf)
	}
}

//随机源（产生排序的原始数据）
func RandomSource(count int) <-chan int {
	out := make(chan int)
	go func() {
		for i := 0; i < count; i++ {
			out <- rand.Int()
		}
		close(out)
	}()
	return out
}

//多路合并
func MergeN(inputs ...<-chan int) <-chan int {
	if len(inputs) == 1 {
		return inputs[0]
	}
	m := len(inputs) / 2
	return Merge(MergeN(inputs[:m]...), MergeN(inputs[m:]...))
}
