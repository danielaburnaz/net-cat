package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"nc/connection"
)

func main() {
	address, err := connection.ParseArgs(os.Args)
	if err != nil {
		log.Println("Usage: go run . $IP [$port]")
		return
	}

	conn, err := connection.NewChatConnection(address, onMessage, onUsersUpdated)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	stdin := bufio.NewScanner(os.Stdin)
	for stdin.Scan() {
		switch stdin.Text() {
		case "/users":
			fmt.Print("Users: ")
			fmt.Println(formatUserList(conn.Users))
		case "/quit":
			return
		default:
			err = conn.SendMessage(stdin.Text())
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func onMessage(msg string) {
	fmt.Println(msg)
}

// Update user list
func onUsersUpdated(users []string) {}

// For /users format the user list
func formatUserList(userList []string) string {
	var users string
	for i, user := range userList {
		if i > 0 {
			users += ", "
		}
		users += user
	}
	return users
}
