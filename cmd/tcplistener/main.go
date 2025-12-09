package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	// readAndPrintFileInChunks("messages.txt")
	readAndPrintLinesFromTcpConn()
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	c := make(chan string, 1)
	reader := bufio.NewReader(f)
	str := ""
	go func() {
		defer close(c)
		for {
			data := make([]byte, 8)
			n, err := reader.Read(data)
			if n == 0 {
				if err == io.EOF {
					f.Close()
					return
				}
				if err != nil {
					fmt.Println("Error scanning" + err.Error())
					return
				}
			}
			data = data[:n]
			if i := bytes.IndexByte(data, '\n'); i != -1 {
				str += string(data[:i])
				data = data[i+1:]
				if len(str) > 0 {
					c <- str
				}
				str = ""
			}
			str += string(data)
		}
	}()
	return c
}

func readAndPrintLinesFromTcpConn() {
	addr := "127.0.0.1:42069"
	listner, err := net.Listen("tcp", addr)

	if err != nil {
		fmt.Println("Error listening on network: " + err.Error())
		os.Exit(-1)
	}
	defer listner.Close()

	fmt.Println("Listening connections on " + addr)

	for {
		conn, err := listner.Accept()
		if err != nil {
			fmt.Println("Error accepting a connection: " + err.Error())
			log.Fatal("error", err)
		}
		for s := range getLinesChannel(conn) {
			fmt.Printf("read: %s\n", s)
		}
	}
	fmt.Println("Connection has been closed")
}

func readAndPrintFileInChunks(filepath string) {
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	c := getLinesChannel(file)
	for s := range c {
		fmt.Printf("read: %s\n", s)
	}
}
