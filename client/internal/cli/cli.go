package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/google/shlex"
	"github.com/z3vxo/kronos/internal/httpclient"
)

type CLI struct {
	http          *httpclient.Client
	ui            *UI
	ClientInUse   string
	dispatchTable map[string]HandlerFunc
}

type HandlerFunc func(args []string)

func NewCli() (*CLI, error) {
	h, err := httpclient.NewClient()
	if err != nil {
		return nil, err
	}

	rl, err := NewUI()
	if err != nil {
		return nil, err
	}

	c := &CLI{
		http: h,
		ui:   rl,
	}
	go c.ui.Run()
	c.SetupDispatchTable()

	go h.ConnectToSSE()

	return c, nil
}

func (c *CLI) SetupDispatchTable() {
	c.dispatchTable = map[string]HandlerFunc{
		"list":      c.ListAgents,
		"use":       c.ResolveAgent,
		"back":      c.Back,
		"listeners": c.ListListners,
	}
}

func (c *CLI) Close() {
	c.ui.rl.Close()
}

func (c *CLI) Dispatch(cmd []string) {
	fn, ok := c.dispatchTable[cmd[0]]
	if !ok {
		c.ui.Send(fmt.Sprintf("[!] Unknown command: %s", cmd[0]))
		return
	}
	go fn(cmd[1:])
}

func (c *CLI) Run() {
	defer c.Close()
	for {
		input, err := c.ui.rl.Readline()
		if err == io.EOF || err == readline.ErrInterrupt {
			break
		}
		if err != nil {
			c.ui.Send(fmt.Sprintf("error: %v", err))
			break
		}
		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		cmd, err := shlex.Split(input)
		if err != nil {
			c.ui.Send("[!] Failed parsing input")
			continue
		}

		if strings.ToLower(cmd[0]) == "exit" {
			c.Close()
			os.Exit(0)
		}

		c.ui.rl.SaveHistory(input)
		c.Dispatch(cmd)
	}
}
