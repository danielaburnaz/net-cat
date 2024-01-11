package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

const (
	ENTER_NAME = "[ENTER YOUR NAME]: "
)

var connections []net.Conn

func main() {
	arg := os.Args

	port := ""
	switch len(arg) {
	case 1:
		port = "8989"
	case 2:
		port = arg[1]
	default:
		log.Println("[USAGE]: ./TCPChat $port")
		return
	}

	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Listening on the port :%s\n", port)

	ch := make(chan string)
	go channelMessages(ch)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
		}
		go handleConnection(conn, ch)
		connections = append(connections, conn)
	}
}

func channelMessages(ch chan string) {
	for msg := range ch {
		for _, conn := range connections {
			conn.Write([]byte(msg))
		}
	}
}

func handleConnection(conn net.Conn, ch chan string) {
	scanner := bufio.NewScanner(conn)
	username := name(conn, scanner)

	ch <- fmt.Sprintf("%s joined\n", username)

	for scanner.Scan() {
		if scanner.Text() != "" {
			// fmt.Print(format(username, scanner.Text()))

			ch <- format(username, scanner.Text())
		}
	}
	ch <- fmt.Sprintf("%s left\n", username)

}

func format(username string, text string) string {
	time := time.Now().Format(time.DateTime)
	return fmt.Sprintf("[%s][%s]: %s\n", time, username, text)
}

func name(conn net.Conn, scanner *bufio.Scanner) string {
	conn.Write([]byte(ENTER_NAME))

	for scanner.Scan() {
		if strings.TrimSpace(scanner.Text()) != "" {
			return scanner.Text()
		}
		conn.Write([]byte(ENTER_NAME))

	}
	return "unknown"
}
