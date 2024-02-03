package main

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

var (
	chatMessages []string
	userList     []string
)

func main() {
	userList = append(userList, "user1", "user2", "user3", "user4", "user5", "user6", "user7", "user8")
	chatMessages = append(chatMessages, "Hello", "World", "How", "Are", "You", "Today", "I", "Am", "Fine", "Thank", "You", "For", "Asking", "Goodbye")
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Fatalf("Failed to create GUI: %v", err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := keybindings(g); err != nil {
		log.Fatalf("Failed to set keybindings: %v", err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Fatalf("Failed main loop: %v", err)
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
		for _, user := range userList {
			fmt.Fprintln(v, user)
		}
	}

	if v, err := g.SetView("chat", 21, 0, maxX-1, maxY-5); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Chat"
		v.Wrap = true
		// v.Autoscroll = true
		v.Clear()
		for _, msg := range chatMessages {
			fmt.Fprintln(v, msg)
		}
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

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	if err := g.SetKeybinding("input", gocui.KeyEnter, gocui.ModNone, sendMessage); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, scroll(true)); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, scroll(false)); err != nil {
		return err
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func sendMessage(g *gocui.Gui, v *gocui.View) error {
	inputView, _ := g.View("input")
	message := inputView.Buffer()

	// Process the message (e.g., send it to a chat server)
	chatMessages = append(chatMessages, message)

	// Clear the input field
	inputView.Clear()
	inputView.SetCursor(0, 0)

	// Update the chat view
	chatView, _ := g.View("chat")
	fmt.Fprint(chatView, message)

	return nil
}

func inputEditor(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	if key != gocui.KeyEnter {
		gocui.DefaultEditor.Edit(v, key, ch, mod)
	}
}

func scroll(up bool) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		v.Autoscroll = false
		chatView, err := g.View("chat")
		if err != nil {
			return err
		}

		// _, y := v.Size()
		ox, oy := chatView.Origin()
		if up {
			// v.Autoscroll = false
			if oy > 0 {
				chatView.SetOrigin(ox, oy-1)
			}
		} else {
			// v.Autoscroll = false
			chatView.SetOrigin(ox, oy+1)
		}
		return nil
	}
}
