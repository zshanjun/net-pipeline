package pipeline

import (
	"bufio"
	"net"
)

//网络写入
func NetworkSink(addr string, in <-chan int) {
	//开启一个服务器，将准备好的数据写入到服务器，等待客户端连接读取
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	go func() {
		defer listener.Close()

		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		defer conn.Close()

		writer := bufio.NewWriter(conn)
		//使用bufio后，需要flush才会将数据写入
		defer writer.Flush()
		WriterSink(writer, in)
	}()
}

//网络来源
func NetworkSource(addr string) <-chan int {
	out := make(chan int)
	go func() {
		//启动一个客户端读取数据
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			panic(err)
		}

		r := ReaderSource(bufio.NewReader(conn), -1)
		for v := range r {
			out <- v
		}
		close(out)
	}()
	return out
}
