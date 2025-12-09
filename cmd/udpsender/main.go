package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	runUDPSender()
}

func runUDPSender() {
	addr := "localhost:42069"

	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Fatal("Error resolving UDP addr", err)
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Fatal("Error dialing UDP connection ", err)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf(">\n")
		str, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from stdin: " + err.Error())
			continue
		}
		_, err = conn.Write([]byte(str))
		if err != nil {
			fmt.Println("Error writing to connection: " + err.Error())
			continue
		}

	}
}
