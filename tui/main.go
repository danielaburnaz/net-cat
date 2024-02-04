package main

import (
	"fmt"
	"log"
	"os"

	"nc/connection"

	"github.com/jroimartin/gocui"
)

func main() {
	address, err := connection.ParseArgs(os.Args)
	if err != nil {
		log.Println("Usage: go run . $IP [$port]")
		return
	}
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Fatalf("Failed to create GUI: %v", err)
	}
	defer g.Close()

	conn, err := connection.NewChatConnection(address, onMessage(g), onUsersUpdated(g))
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	g.SetManagerFunc(layout)

	if err := keybindings(g, conn); err != nil {
		log.Fatalf("Failed to set keybindings: %v", err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Fatalf("Failed main loop: %v", err)
	}
}

func onMessage(g *gocui.Gui) func(string) {
	return func(msg string) {
		chatView, _ := g.View("chat")
		fmt.Fprintln(chatView, msg)
		g.Update(func(g *gocui.Gui) error {
			return nil
		})
	}
}

func onUsersUpdated(g *gocui.Gui) func([]string) {
	return func(users []string) {
		userListView, _ := g.View("userList")
		userListView.Clear()
		for _, user := range users {
			fmt.Fprintln(userListView, user)
		}
		g.Update(func(g *gocui.Gui) error {
			return nil
		})
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("userList", 0, 0, 20, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Users"
		v.Wrap = true
		v.Autoscroll = true
		v.Clear()
	}

	if v, err := g.SetView("chat", 21, 0, maxX-1, maxY-5); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Chat"
		v.Wrap = true
		v.Autoscroll = true
		v.Clear()
	}

	if v, err := g.SetView("input", 21, maxY-4, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Input"
		v.Editable = true
		v.Wrap = true
		v.Autoscroll = true
		v.Editor = gocui.EditorFunc(inputEditor)
		if _, err := g.SetCurrentView("input"); err != nil {
			return err
		}
	}

	return nil
}

func keybindings(g *gocui.Gui, conn *connection.ChatConnection) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	if err := g.SetKeybinding("input", gocui.KeyEnter, gocui.ModNone, sendMessage(conn)); err != nil {
		return err
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func sendMessage(conn *connection.ChatConnection) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		inputView, _ := g.View("input")
		message := inputView.Buffer()

		conn.SendMessage(message)

		// Clear the input field
		inputView.Clear()
		inputView.SetCursor(0, 0)

		return nil
	}
}

func inputEditor(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	if key != gocui.KeyEnter {
		gocui.DefaultEditor.Edit(v, key, ch, mod)
	}
}
