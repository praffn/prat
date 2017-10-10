package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"prat/prat"
	"strings"

	tui "github.com/marcusolsson/tui-go"
)

var name = flag.String("name", "anon", "your name")
var host = flag.String("host", "localhost", "the host to connect to")
var port = flag.Int("port", prat.DefaultPort, "the port to connect to")
var server = flag.Bool("server", false, "starts a server instead of a client")
var logFile = flag.String("log", "prat.log", "file to output logging information to")

func main() {
	flag.Parse()
	if *server {
		file, err := os.Create(*logFile)
		if err != nil {
			panic(err)
		}
		logger := log.New(file, "", 0)
		server := prat.NewServerWithLogger(logger)
		server.Start(*port)
	} else {
		CreateClient(*name, *host, *port)
	}
}

func CreateClient(name, host string, port int) {
	fmt.Println(name, host)
	history := tui.NewVBox()
	history.SetBorder(true)
	history.Append(tui.NewSpacer())
	input := tui.NewEntry()
	input.SetFocused(true)
	input.SetSizePolicy(tui.Expanding, tui.Maximum)
	inputBox := tui.NewHBox(input)
	inputBox.SetBorder(true)
	inputBox.SetSizePolicy(tui.Expanding, tui.Maximum)

	chat := tui.NewVBox(history, inputBox)
	chat.SetSizePolicy(tui.Expanding, tui.Expanding)
	ui := tui.New(chat)
	ui.SetKeybinding("Esc", func() { ui.Quit() })

	address := fmt.Sprintf("%s:%d", host, port)
	client := prat.NewClient(address, name, func(msg prat.Message) {
		history.Append(tui.NewHBox(
			tui.NewLabel(msg.Time.Format("15:04")),
			tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("<%s>", msg.Author))),
			tui.NewLabel(msg.Body),
			tui.NewSpacer(),
		))
		ui.Update(func() {})
	})

	input.OnSubmit(func(e *tui.Entry) {
		text := e.Text()
		if len(text) == 0 {
			return
		}
		if text[0] == '/' {
			split := strings.Split(text, " ")
			if len(split) < 2 {
				input.SetText("")
				return
			}
			switch split[0] {
			case "/setname":
				client.Name = split[1]
			}
			input.SetText("")
			return
		}

		client.SendMessage(text)
		input.SetText("")
	})

	if err := ui.Run(); err != nil {
		panic(err)
	}
}
