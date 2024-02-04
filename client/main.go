package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"

	"github.com/jroimartin/gocui"
)

var userList []string
var chatMessages []string
var conn net.Conn

func main() {
	arg := os.Args

	port := ""

	switch len(arg) {
	case 2:
		port = "8989"
	case 3:
		port = arg[2]
	default:
		log.Println("Usage: go run . $IP [$port]")
		return
	}

	ip := arg[1]

	conn, err := net.Dial("tcp", ip+":"+port)
	if err != nil {
		log.Println(err)
		return
	}

	// defer conn.Close()
	g := gui()

	// stdin := bufio.NewScanner(os.Stdin)
	read := bufio.NewScanner(conn)
	read.Split(bufio.ScanLines)
	go reader(g, read)
	// for stdin.Scan() {
	// 	// _, err = conn.Write(append(stdin.Bytes(), byte('\n')))

	// 	// fmt.Println(userList)
	// 	if err != nil {
	// 		return
	// 	}
	// }
}

func reader(g *gocui.Gui, read *bufio.Scanner) {
	for read.Scan() {
		// Match user joined/left messages and update user list

		j := regexp.MustCompile(`^(.+) has joined our chat$`)
		joined := j.MatchString(read.Text())

		l := regexp.MustCompile(`^(.+) has left our chat$`)
		left := l.MatchString(read.Text())

		if joined {
			username := j.FindStringSubmatch(read.Text())
			userList = append(userList, username[1])
			g.Update(func(g *gocui.Gui) error {
				v, err := g.View("userList")
				if err != nil {
					panic(err)
				}
				v.Clear()
				for _, user := range userList {
					fmt.Fprintln(v, user)
				}

				return nil
			})
		}
		if left {
			username := l.FindStringSubmatch(read.Text())
			removeUsername(username[1])
			g.Update(func(g *gocui.Gui) error {
				v, err := g.View("userList")
				if err != nil {
					panic(err)
				}
				v.Clear()
				for _, user := range userList {
					fmt.Fprintln(v, user)
				}

				return nil
			})
		}
		chatMessages = append(chatMessages, fmt.Sprint(read.Text()))
		g.Update(func(g *gocui.Gui) error {
			v, err := g.View("chat")
			if err != nil {
				panic(err)
			}
			fmt.Fprintln(v, read.Text())

			return nil
		})
	}
}

func removeUsername(username string) {
	for i, u := range userList {
		if u == username {
			userList = append(userList[:i], userList[i+1:]...)
			return
		}
	}
}
