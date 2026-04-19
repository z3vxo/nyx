package cli

import (
	"fmt"
	"time"

	"github.com/chzyer/readline"
)

const (
	dim = "\001\033[2m\033[4m\002"
	rst = "\001\033[0m\002"
)

func (u *UI) SetPrompt(agent string) {
	t := time.Now().Format("15:04:05")

	if agent == "" {
		u.rl.SetPrompt(fmt.Sprintf("[%s] %skronos%s $> ", t, dim, rst))
	} else {
		u.rl.SetPrompt(fmt.Sprintf("[%s] %skronos%s (\001\033[33m\002%s%s) $> ", t, dim, rst, agent, rst))
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
