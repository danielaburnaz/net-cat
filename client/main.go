package main

import (
	"bufio"
	"log"
	"net"
	"os"
)

func main() {
	arg := os.Args
	
	port := ""

	switch len(arg){
	case 2:		
		port = "8989"
	case 3:
		port = arg[2]
	default:
		log.Println("Usage: go run . $IP [$port]")
		return
	}
	
	ip := arg[1]

	conn, err := net.Dial("tcp", ip + ":" +port)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Listening to port: :" + port)
	stdin := bufio.NewScanner(os.Stdin)
	read := bufio.NewScanner(conn)
	go reader(read)
	for stdin.Scan() {
		_, err = conn.Write(append(stdin.Bytes(), byte('\n')))
		if err != nil {
			return
		}
	}
}

func reader(read *bufio.Scanner) {
	for read.Scan() {
		log.Println(read.Text())
	}
}
