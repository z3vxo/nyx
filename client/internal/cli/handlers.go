package cli

import (
	"fmt"
	"time"
)

func (c *CLI) ListAgents(args []string) {
	var A Agents
	if err := c.http.DoGet("GET", "ts/rest/agents/list", &A); err != nil {
		fmt.Println("[!] Failed listing agents:", err)
		return
	}

	if A.Total == 0 {
		fmt.Println("[*] No agents connected")
		return
	}

	fmt.Printf("%-12s %-16s %-20s %-16s %-16s %-8s %-6s %s\n",
		"CODENAME", "USER", "HOSTNAME", "EXT IP", "INT IP", "ELEV", "PID", "LAST SEEN")
	fmt.Println("-----------------------------------------------------------------------------------------------------")
	for _, a := range A.Agent {
		elev := "no"
		if a.IsElevated {
			elev = "yes"
		}
		last := time.Unix(a.LastSeen, 0).Format("2006-01-02 15:04:05")
		fmt.Printf("%-12s %-16s %-20s %-16s %-16s %-8s %-6d %s\n",
			a.CodeName, a.Username, a.Hostname, a.Ex_ip, a.In_ip, elev, a.Pid, last)
	}
}
