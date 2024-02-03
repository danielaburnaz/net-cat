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
	BT         = "`"
	TUX        = `
Welcome to TCP-Chat!
         _nnnn_
        dGGGGMMb
       @p~qp~~qMb
       M|@||@) M|
       @,----.JM|
      JS^\__/  qKL
     dZP        qKRb
    dZP          qKKb
   fZP            SMMb
   HZM            MMMM
   FqM            MMMM
 __| ".        |\dS"qML
 |    ` + BT + `.       | ` + BT + `' \Zq
_)      \.___.,|     .'
\____   )MMMMMP|   .'
     ` + BT + `-'       ` + BT + `--'
`
)

var (
	connections    []net.Conn
	messageHistory []byte
)

func removeConnection(conn net.Conn) {
	for i, c := range connections {
		if c == conn {
			connections = append(connections[:i], connections[i+1:]...)
			return
		}
	}
}

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
		switch {
		case len(connections) < 11:
			connections = append(connections, conn)
			go handleConnection(conn, ch)
		default:
			conn.Write([]byte("Max connection reached\n"))
			ch <- "User tried to join but max connection reached\n"
			conn.Close()
			removeConnection(conn)
		}
	}
}

func channelMessages(ch chan string) {
	for msg := range ch {
		for _, conn := range connections {
			conn.Write([]byte(msg))
			messageHistory = append(messageHistory, []byte(msg)...)
		}
	}
}

func handleConnection(conn net.Conn, ch chan string) {
	defer removeConnection(conn)
	scanner := bufio.NewScanner(conn)
	conn.Write([]byte(TUX))

	username := name(conn, scanner)

	conn.Write(messageHistory)
	ch <- fmt.Sprintf("%s joined\n", username)

	for scanner.Scan() {
		if scanner.Text() != "" {
			log.Print(format(username, scanner.Text()))

			ch <- format(username, scanner.Text())
		}
	}
	conn.Close()
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
