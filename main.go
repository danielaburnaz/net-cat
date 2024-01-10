package main

import (
	"bufio"
	"log"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", ":8989")
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
	log.Println("User joined")
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		log.Println(scanner.Text())
		conn.Write([]byte(reverseString(scanner.Text())))
		conn.Write([]byte("\n"))
	}
	log.Println("User left")
}

func reverseString(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}
