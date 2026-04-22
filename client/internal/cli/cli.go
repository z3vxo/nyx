package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/google/shlex"
	"github.com/z3vxo/kronos/internal/httpclient"
	"github.com/z3vxo/kronos/internal/ui"
)

type CLI struct {
	http          *httpclient.Client
	ui            *ui.UI
	ClientInUse   string
	dispatchTable map[string]HandlerFunc
}

type HandlerFunc func(args []string)

func NewCli() (*CLI, error) {
	rl, err := ui.NewUI()
	if err != nil {
		return nil, err
	}
	h, err := httpclient.NewClient(rl)
	if err != nil {
		return nil, err
	}

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
		"listeners": c.ParseListenerCmd,
		"info":      c.ListAgentInfo,
		"help":      c.Help,
	}
}

func (c *CLI) Close() {
	c.ui.Rl.Close()
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
		input, err := c.ui.Rl.Readline()
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
			c.ui.Send(ui.BAD.Sprint("Failed parsing input"))
			continue
		}

		if strings.ToLower(cmd[0]) == "exit" {
			c.Close()
			os.Exit(0)
		}
		c.ui.Rl.SaveHistory(input)

		c.Dispatch(cmd)
	}
}

func (c *CLI) Help(args []string) {
	c.ui.Send("\n")
	c.ui.Send("\033[1;35m  ██╗  ██╗██████╗  ██████╗ ███╗   ██╗ ██████╗ ███████╗\033[0m")
	c.ui.Send("\033[1;35m  ██║ ██╔╝██╔══██╗██╔═══██╗████╗  ██║██╔═══██╗██╔════╝\033[0m")
	c.ui.Send("\033[1;35m  █████╔╝ ██████╔╝██║   ██║██╔██╗ ██║██║   ██║███████╗\033[0m")
	c.ui.Send("\033[1;35m  ██╔═██╗ ██╔══██╗██║   ██║██║╚██╗██║██║   ██║╚════██║\033[0m")
	c.ui.Send("\033[1;35m  ██║  ██╗██║  ██║╚██████╔╝██║ ╚████║╚██████╔╝███████║\033[0m")
	c.ui.Send("\033[1;35m  ╚═╝  ╚═╝╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚═══╝ ╚═════╝ ╚══════╝\033[0m")
	c.ui.Send("")
	c.ui.Send("\033[1;37m  AGENTS\033[0m")
	c.ui.Send(fmt.Sprintf("  \033[1;36m%-40s\033[0m %s", "list", "list all connected agents"))
	c.ui.Send(fmt.Sprintf("  \033[1;36m%-40s\033[0m %s", "use <codename>", "interact with an agent"))
	c.ui.Send(fmt.Sprintf("  \033[1;36m%-40s\033[0m %s", "info", "detailed info on current agent"))
	c.ui.Send(fmt.Sprintf("  \033[1;36m%-40s\033[0m %s", "back", "stop using current agent"))
	c.ui.Send("")
	c.ui.Send("\033[1;37m  LISTENERS\033[0m")
	c.ui.Send(fmt.Sprintf("  \033[1;36m%-40s\033[0m %s", "listeners", "list active listeners"))
	c.ui.Send(fmt.Sprintf("  \033[1;36m%-40s\033[0m %s", "listeners new -h <host> -p <port> -t <proto>", "start listener (proto: http|https)"))
	c.ui.Send(fmt.Sprintf("  \033[1;36m%-40s\033[0m %s", "listeners start <name>", "start a listener"))
	c.ui.Send(fmt.Sprintf("  \033[1;36m%-40s\033[0m %s", "listeners stop <name>", "stop a listener"))
	c.ui.Send("")
	c.ui.Send("\033[1;37m  COMMANDS\033[0m")
	c.ui.Send(fmt.Sprintf("  \033[1;36m%-40s\033[0m %s", "getprivs", "get current users privileges"))
	c.ui.Send(fmt.Sprintf("  \033[1;36m%-40s\033[0m %s", "cd <dir>", "change directory"))
	c.ui.Send(fmt.Sprintf("  \033[1;36m%-40s\033[0m %s", "ls <dir>", "list a directory"))
	c.ui.Send(fmt.Sprintf("  \033[1;36m%-40s\033[0m %s", "cat <file>", "read a file"))
	c.ui.Send(fmt.Sprintf("  \033[1;36m%-40s\033[0m %s", "whoami", "list current users identity"))
	c.ui.Send(fmt.Sprintf("  \033[1;36m%-40s\033[0m %s", "WIP MORE CMDS COMING"))

}
