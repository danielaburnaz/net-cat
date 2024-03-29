// common vars between the client and UI
package connection

import (
	"bufio"
	"net"
	"regexp"
)

var (
	JOINED = regexp.MustCompile(`^(.+) has joined our chat$`)
	LEFT   = regexp.MustCompile(`^(.+) has left our chat$`)
)

type ChatConnection struct {
	conn  net.Conn
	Users []string
}

// Check os.Args specific port
func ParseArgs(args []string) (string, error) {
	if len(args) < 2 {
		return "", nil
	}
	ip := ""
	port := "8989"
	if len(args) > 2 {
		port = args[2]
	}
	return ip + ":" + port, nil
}

// Update Client UI to display new message/user
func NewChatConnection(address string, onMessage func(string), onUsersUpdated func([]string)) (*ChatConnection, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	chatConnection := &ChatConnection{conn: conn}
	go chatConnection.readMessages(onMessage, onUsersUpdated)

	return chatConnection, nil
}

// Check for new users/user left and read new messages
func (c *ChatConnection) readMessages(onMessage func(string), onUsersUpdated func([]string)) {
	scanner := bufio.NewScanner(c.conn)
	for scanner.Scan() {
		line := scanner.Text()

		if match := JOINED.FindStringSubmatch(line); len(match) > 0 {
			c.Users = append(c.Users, match[1])
			onUsersUpdated(c.Users)
		}
		if match := LEFT.FindStringSubmatch(line); len(match) > 0 {
			c.Users = removeFromList(c.Users, match[1])
			onUsersUpdated(c.Users)
		}
		onMessage(line)
	}
}

// Close connection
func (c *ChatConnection) Close() {
	c.conn.Close()
}

// Send message
func (c *ChatConnection) SendMessage(message string) error {
	_, err := c.conn.Write([]byte(message + "\n"))
	return err
}

// Remove user from the list once they leave
func removeFromList(list []string, item string) []string {
	for i, u := range list {
		if u == item {
			return append(list[:i], list[i+1:]...)
		}
	}
	return list
}
