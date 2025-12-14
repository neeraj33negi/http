package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/neeraj33negi/http/internal/request"
)

func main() {
	readAndPrintLinesFromTcpConn()
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
		r, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Println("Error reading request " + err.Error())
			log.Fatal("error", err)
		}
		fmt.Printf("Method: %s\nHttpVersion: %s\nTarget: %s\n",
			r.RequestLine.Method,
			r.RequestLine.HttpVersion,
			r.RequestLine.RequestTarget,
		)
	}
}
