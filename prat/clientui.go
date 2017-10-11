package prat

import (
	"fmt"
	"strings"

	tui "github.com/marcusolsson/tui-go"
)

type ClientUI struct {
	history *tui.Box
	input   *tui.Entry
	ui      tui.UI
	client  *Client
}

func NewClientUI(client *Client) *ClientUI {
	history := tui.NewVBox(tui.NewSpacer())
	history.SetBorder(true)

	input := tui.NewEntry()
	input.SetFocused(true)
	input.SetSizePolicy(tui.Expanding, tui.Maximum)
	inputBox := tui.NewHBox(input)
	inputBox.SetBorder(true)
	inputBox.SetSizePolicy(tui.Expanding, tui.Maximum)

	root := tui.NewVBox(history, inputBox)
	root.SetSizePolicy(tui.Expanding, tui.Expanding)
	ui := tui.New(root)
	ui.SetKeybinding("Esc", func() { ui.Quit() })

	cui := ClientUI{
		history: history,
		input:   input,
		ui:      ui,
		client:  client,
	}

	client.AddMessageHandler(cui.AddMessageToHistory)
	input.OnSubmit(cui.OnMessageSubmit)

	return &cui
}

func (cui *ClientUI) Run() error {
	return cui.ui.Run()
}

func (cui *ClientUI) AddMessageToHistory(message Message) {
	cui.history.Append(tui.NewHBox(
		tui.NewLabel(message.Time.Format("15:04")),
		tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("<%s>", message.Author))),
		tui.NewLabel(message.Body),
		tui.NewSpacer(),
	))
	cui.ui.Update(func() {})
}

func (cui *ClientUI) OnMessageSubmit(e *tui.Entry) {
	// defer resetting textfield
	defer e.SetText("")
	// get text from entry
	text := e.Text()
	if len(text) == 0 {
		// if empty do nothing
		return
	}
	if text[0] == '/' {
		// if first character is a slash, the message
		// is a command. We split the message into head
		// and rest (e.g. head: "/setname", rest: ["phin"])
		// we also remove the first slash
		args := strings.Split(text, " ")
		command := args[0][1:]
		rest := args[1:]
		cui.ParseCommand(command, rest)
	} else {
		// otherwise, send as message
		cui.client.SendMessage(e.Text())
	}
}

func (cui *ClientUI) ParseCommand(command string, rest []string) {
	switch command {
	case "setname":
		if len(rest) > 0 {
			cui.client.Name = rest[0]
		}
	case "exit":
		cui.ui.Quit()
	case "help":
		cui.PrintHelp()
	}
}

func (cui *ClientUI) PrintHelp() {
	cui.history.Append(
		tui.NewPadder(1, 1, tui.NewVBox(
			tui.NewLabel("Commands:"),
			tui.NewHBox(
				tui.NewLabel("/setname <name>"),
				tui.NewPadder(1, 0, tui.NewLabel("Sets your name")),
			),
			tui.NewHBox(
				tui.NewLabel("/help"),
				tui.NewPadder(1, 0, tui.NewLabel("Prints help")),
			),
			tui.NewHBox(
				tui.NewLabel("/exit"),
				tui.NewPadder(1, 0, tui.NewLabel("Terminates the session")),
			),
		)),
	)
}
