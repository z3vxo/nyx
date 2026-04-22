package cli

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/z3vxo/kronos/internal/ui"
)

func relativeTime(unix int64) string {
	since := time.Since(time.Unix(unix, 0))
	switch {
	case since < time.Minute:
		return fmt.Sprintf("%ds", int(since.Seconds()))
	case since < time.Hour:
		return fmt.Sprintf("%dm", int(since.Minutes()))
	case since < 24*time.Hour:
		return fmt.Sprintf("%dh", int(since.Hours()))
	default:
		return fmt.Sprintf("%dd", int(since.Hours()/24))
	}
}

func (c *CLI) ListAgents(args []string) {
	var A Agents
	if err := c.http.DoGet("ts/rest/agents/list", &A); err != nil {
		c.ui.Send(ui.WARN.Sprintf("Failed listing agents: %v", err))
		return
	}

	if A.Total == 0 {
		c.ui.Send(ui.INFO.Sprint("No agents connected"))
		return
	}

	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "--------\t----\t--------\t------\t------\t----\t---\t---------\t--------)")
	fmt.Fprintln(w, "CODENAME\tUSER\tHOSTNAME\tEXT IP\tINT IP\tELEV\tPID\tLAST SEEN\tREG DATE")
	fmt.Fprintln(w, "--------\t----\t--------\t------\t------\t----\t---\t---------\t--------")
	for _, a := range A.Agent {
		elev := "no"
		if a.IsElevated {
			elev = "yes"
		}
		last := relativeTime(a.LastSeen)
		reg := time.Unix(a.RegDate, 0).Format("2006-01-02 15:04:05")
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%d\t%s\t%s\n",
			a.CodeName, a.Username, a.Hostname, a.Ex_ip, a.In_ip, elev, a.Pid, last, reg)
	}

	w.Flush()
	c.ui.Send(buf.String())
}

func (c *CLI) ListAgentInfo(args []string) {
	if c.ui.InUse == "" {
		c.ui.Send(ui.BAD.Sprint("Must be using agent!"))
		return
	}
	var a AgentInfoResp

	if err := c.http.DoGet(fmt.Sprintf("ts/rest/agents/info/%s", c.ui.InUse), &a); err != nil {
		c.ui.Send(ui.WARN.Sprintf("Failed listing info: %v", err))
		return
	}

	elev := "no"
	if a.IsElevated {
		elev = "yes"
	}
	arch := "x86"
	if a.Arch == 1 {
		arch = "x64"
	}
	last := relativeTime(a.LastCheckin)
	reg := time.Unix(a.RegisterTime, 0).Format("2006-01-02 15:04:05")

	rows := [][2]string{
		{"CODENAME", c.ui.InUse},
		{"USER", a.User},
		{"HOSTNAME", a.Host},
		{"INT IP", a.InternalIP},
		{"EXT IP", a.ExternalIP},
		{"PID / PPID", fmt.Sprintf("%d / %d", a.Pid, a.PPid)},
		{"ARCH", arch},
		{"ELEVATED", elev},
		{"OS", a.WinVer},
		{"PROCESS", a.ProcPath},
		{"LAST SEEN", last},
		{"REGISTERED", reg},
	}

	col0, col1 := 0, 0
	for _, r := range rows {
		if len(r[0]) > col0 {
			col0 = len(r[0])
		}
		if len(r[1]) > col1 {
			col1 = len(r[1])
		}
	}

	sep := fmt.Sprintf("+-%-*s-+-%-*s-+", col0, strings.Repeat("-", col0), col1, strings.Repeat("-", col1))

	var buf bytes.Buffer
	buf.WriteString(sep + "\n")
	for _, r := range rows {
		fmt.Fprintf(&buf, "| %-*s | %-*s |\n", col0, r[0], col1, r[1])
	}
	buf.WriteString(sep + "\n")
	c.ui.Send(buf.String())
}

func (c *CLI) Back(args []string) {
	c.ui.Send(ui.INFO.Sprintf("Not using %s", c.ClientInUse))
	c.ClientInUse = ""
	c.ui.InUse = ""
	c.ui.SetPrompt("")
}

func (c *CLI) ResolveAgent(args []string) {
	if len(args) < 1 || args[0] == "" {
		c.ui.Send(ui.BAD.Sprint("Error: must choose agent"))
		return
	}

	var r ResolveResp
	e := fmt.Sprintf("ts/rest/agents/resolve/%s", args[0])

	if err := c.http.DoGet(e, &r); err != nil {
		c.ui.Send(ui.WARN.Sprintf("Failed resolving agent: %v", err))
		return
	}

	if r.Guid == "" {
		c.ui.Send(ui.BAD.Sprint("Server Did not return a guid!"))
		return
	}

	c.ClientInUse = r.Guid
	c.ui.InUse = args[0]
	c.ui.SetPrompt(args[0])
	c.ui.Send(ui.GOOD.Sprintf("Using %s", c.ClientInUse))
	return
}
