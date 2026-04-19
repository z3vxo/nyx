package cli

import (
	"bytes"
	"fmt"
	"text/tabwriter"
	"time"
)

func (c *CLI) ListAgents(args []string) {
	var A Agents
	if err := c.http.DoGet("ts/rest/agents/list", &A); err != nil {
		c.ui.Send(fmt.Sprintf("[!] Failed listing agents: %v", err))
		return
	}

	if A.Total == 0 {
		c.ui.Send("[*] No agents connected")
		return
	}

	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "--------\t----\t--------\t------\t------\t----\t---\t---------")
	fmt.Fprintln(w, "CODENAME\tUSER\tHOSTNAME\tEXT IP\tINT IP\tELEV\tPID\tLAST SEEN")
	fmt.Fprintln(w, "--------\t----\t--------\t------\t------\t----\t---\t---------")
	for _, a := range A.Agent {
		elev := "no"
		if a.IsElevated {
			elev = "yes"
		}
		last := time.Unix(a.LastSeen, 0).Format("2006-01-02 15:04:05")
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%d\t%s\n",
			a.CodeName, a.Username, a.Hostname, a.Ex_ip, a.In_ip, elev, a.Pid, last)
	}

	w.Flush()
	c.ui.Send(buf.String())
}

func (c *CLI) Back(args []string) {
	c.ui.Send(fmt.Sprintf("[*] Not using %s", c.ClientInUse))
	c.ClientInUse = ""
	c.ui.InUse = ""
	c.ui.SetPrompt("")
}

func (c *CLI) ResolveAgent(args []string) {
	if len(args) < 1 || args[0] == "" {
		c.ui.Send("[!] Error: must choose agent")
		return
	}

	var r ResolveResp
	e := fmt.Sprintf("ts/rest/agents/resolve/%s", args[0])

	if err := c.http.DoGet(e, &r); err != nil {
		c.ui.Send(fmt.Sprintf("[!] Failed resolving agent: %v", err))
		return
	}

	if r.Guid == "" {
		c.ui.Send("[!] Server Did not return a guid!")
		return
	}

	c.ClientInUse = r.Guid
	c.ui.InUse = args[0]
	c.ui.SetPrompt(args[0])
	c.ui.Send(fmt.Sprintf("[*] Using %s", c.ClientInUse))
	return
}

func (c *CLI) ListListners(args []string) {
	var r ListListenersResp

	if err := c.http.DoGet("ts/rest/listeners/list", &r); err != nil {
		c.ui.Send("[!] Failed Listing Listeners!")
		return
	}

	if r.Total == 0 {
		c.ui.Send("[+] No Active Listeners")
		return
	}

	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "--------\t----\t--------\t------")
	fmt.Fprintln(w, "NAME\tPORT\tPROTOCOL\tSTATUS")
	fmt.Fprintln(w, "--------\t----\t--------\t------")

	for _, i := range r.Listeners {
		fmt.Fprintf(w, "%s\t%d\t%s\t%s\n",
			i.Name, i.Port, "HTTPS", "RUNNING")
	}
	w.Flush()
	c.ui.Send(buf.String())
}
