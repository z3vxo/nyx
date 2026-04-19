package cli

import (
	"fmt"
	"time"

	"github.com/chzyer/readline"
)

func (u *UI) SetPrompt(agent string) {
	t := time.Now().Format("15:04:05")
	if agent == "" {
		u.rl.SetPrompt(fmt.Sprintf("\001\033[2m\033[4m\002[%s] kronos\001\033[0m\002 $> ", t))
	} else {
		u.rl.SetPrompt(fmt.Sprintf("\001\033[2m\033[4m\002[%s] kronos\001\033[0m\002(\001\033[33m\002%s\001\033[0m\002) $> ", t,
			agent))
	}
	u.rl.Refresh()
}

func NewUI() (*UI, error) {
	rl, err := readline.NewEx(&readline.Config{
		HistoryFile:     "/tmp/kronos_history",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		return nil, err
	}

	ui := &UI{
		rl:       rl,
		messages: make(chan string, 256),
	}

	ui.SetPrompt("")

	return ui, nil
}

type UI struct {
	messages chan string
	rl       *readline.Instance
	InUse    string
}

func (u *UI) Run() {
	for msg := range u.messages {
		u.rl.Clean()
		fmt.Fprintln(u.rl.Stdout(), msg)
		u.SetPrompt(u.InUse)
		u.rl.Refresh()
	}
}

func (o *UI) Send(msg string) {
	o.messages <- msg
}
