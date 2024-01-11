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
	if len(arg) == 1{
		port = "8989"
	} else if len(arg) != 2{
		log.Println("[USAGE]: ./TCPChat $port")
		return
	} else {
		port = arg[1]
	}

	ln, err := net.Listen("tcp", ":" + port)
	if err != nil {
		panic(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	conn.Write([]byte("[ENTER YOUR NAME]: "))
	scanner.Scan()
	username := scanner.Text()
	log.Println(username + " joined")
	for scanner.Scan() {
		if scanner.Text() != "" {
			log.Println(username + ": " + scanner.Text())
			// conn.Write([]byte(reverseString(scanner.Text())))
			conn.Write([]byte(username + ": " + scanner.Text()))
			conn.Write([]byte("\n"))
		}
	}
	log.Println( username +" left")
}

// func reverseString(s string) string {
// 	r := []rune(s)
// 	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
// 		r[i], r[j] = r[j], r[i]
// 	}
// 	return string(r)
// }
