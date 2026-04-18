package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/google/shlex"
	"github.com/peterh/liner"
	"github.com/z3vxo/kronos/internal/httpclient"
)

type CLI struct {
	http          *httpclient.Client
	liner         *liner.State
	ClientInUse   string
	dispatchTable map[string]HandlerFunc
}

type HandlerFunc func(args []string)

func NewCli() (*CLI, error) {
	h, err := httpclient.NewClient()
	if err != nil {
		return nil, err
	}
	c := &CLI{
		http:  h,
		liner: liner.NewLiner(),
	}
	c.SetupDispatchTable()
	return c, nil
}

func (c *CLI) SetupDispatchTable() {
	c.dispatchTable = map[string]HandlerFunc{
		"list": c.ListAgents,
	}
}

func (c *CLI) Close() {
	c.liner.Close()
}

func (c *CLI) Split(input string) ([]string, error) {
	return shlex.Split(input)
}

func (c *CLI) Dispatch(cmd []string) {
	fn, ok := c.dispatchTable[cmd[0]]
	if !ok {
		fmt.Println("[!] Unknown Command: ", cmd[0])
		return
	}
	fn(cmd[1:])

}

func (c *CLI) Run() {
	defer c.Close()
	for {
		input, err := c.liner.Prompt("kronos $> ")
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("error:", err)
			break
		}
		if input == "" {
			continue
		}

		cmd, err := c.Split(input)
		if err != nil {
			fmt.Println("[!] Failed Parsing input")
		}
		if strings.ToLower(cmd[0]) == "exit" {
			c.Close()
			os.Exit(1)
		}

		c.Dispatch(cmd)

		c.liner.AppendHistory(input)
		fmt.Println(input)
	}
}
