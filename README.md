# Net-Cat: TCP-Chat Server & Client UI in Go

## Overview

Net-Cat is a TCP chat application with both Server and Client UI components written in Go. The server supports a maximum of 10 simultaneous client connections. Upon connecting, users are asked to enter a username, and the chat history is displayed.
The Client UI includes features such as an active user list, a scrollable chat window using arrowkeys, and an input field.

## Usage

### Server

Run the server with the default port (8989):

```bash
go run .
```

Run the server with a custom port:

```bash
go run . $port
```

### Client

Run the client with the server's IP address and port:

```bash
    ~/net-cat$ cd tui
    ~/net-cat/tui$ go run . localhost 8989
```

## UI Package

The client UI is built using the [gocui](https://github.com/jroimartin/gocui) package.

```bash
go get -u github.com/jroimartin/gocui
```

## Features

- **Server:**

  - Listens for incoming connections on the specified port.
  - Manages up to 10 simultaneous client connections.
  - Broadcasts messages to all connected clients through channels.
  - Utilizes goroutines to manage multiple client connections concurrently.

- **Client:**
  - Connects to a server using the specified IP address and port.
  - Prompts the user for a username upon connection.
  - Displays active users, chat history, and an input field for sending messages.


## Contributors

Project done by @spitko and @aburnaz.
