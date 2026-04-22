package ui

import (
	"fmt"
	"time"

	"github.com/chzyer/readline"
)

type tag string

const (
	WARN tag = "\033[1;33m[!]\033[0m "
	INFO tag = "\033[1;36m[*]\033[0m "
	GOOD tag = "\033[1;32m[+]\033[0m "
	BAD  tag = "\033[1;31m[-]\033[0m "
)

func (t tag) Sprint(msg string) string              { return string(t) + msg }
func (t tag) Sprintf(f string, a ...any) string     { return string(t) + fmt.Sprintf(f, a...) }
func (t tag) Sprint_tab(msg string) string          { return string("\t"+t) + msg }
func (t tag) Sprintf_tab(f string, a ...any) string { return string("\t"+t) + fmt.Sprintf(f, a...) }

const (
	dim = "\001\033[2m\033[4m\002"
	rst = "\001\033[0m\002"
)

func (u *UI) SetPrompt(agent string) {
	t := time.Now().Format("15:04:05")

	if agent == "" {
		u.Rl.SetPrompt(fmt.Sprintf("[%s] %skronos%s $> ", t, dim, rst))
	} else {
		u.Rl.SetPrompt(fmt.Sprintf("[%s] %skronos%s (\001\033[33m\002%s%s) $> ", t, dim, rst, agent, rst))
	}
	u.Rl.Refresh()
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
		Rl:       rl,
		messages: make(chan string, 256),
	}

	ui.SetPrompt("")

	return ui, nil
}

type UI struct {
	messages chan string
	Rl       *readline.Instance
	InUse    string
}

func (u *UI) Run() {
	for msg := range u.messages {
		u.Rl.Clean()
		fmt.Fprintln(u.Rl.Stdout(), msg)
		u.SetPrompt(u.InUse)
		u.Rl.Refresh()
	}
}

func (o *UI) Send(msg string) {
	o.messages <- msg
}
