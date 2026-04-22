package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/tabwriter"
	"time"

	"github.com/spf13/pflag"
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
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "--------\t----\t--------\t---------\t------\t------\t----\t---\t--------\t------")
	fmt.Fprintln(w, "CODENAME\tUSER\tHOSTNAME\tPROC PATH\tEXT IP\tINT IP\tELEV\tPID\tLAST SEEN\tREG AT")
	fmt.Fprintln(w, "--------\t----\t--------\t---------\t------\t------\t----\t---\t--------\t------")
	elev := "no"
	if a.IsElevated {
		elev = "yes"
	}
	last := relativeTime(int64(a.LastCheckin))
	reg := time.Unix(int64(a.RegisterTime), 0).Format("2006-01-02 15:04:05")
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%d\t%s\t%s\n",
		c.ui.InUse, a.User, a.Host, a.ProcPath, a.ExternalIP, a.InternalIP, elev, a.Pid, last, reg)

	w.Flush()
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

func (c *CLI) ParseListenerCmd(args []string) {
	if len(args) == 0 {
		c.ListListners()
		return
	}

	switch args[0] {
	case "stop":
		if args[1] == "" {
			c.ui.Send(ui.BAD.Sprint("Must provide listener name"))
		}
		c.StopListener(args[1])
	case "delete":
		if args[1] == "" {
			c.ui.Send(ui.BAD.Sprint("Must provide listener name"))
		}
		c.DeleteListeners(args[1])
	case "start":
		if args[1] == "" {
			c.ui.Send(ui.BAD.Sprint("Must provide listener name"))
		}
		c.StartListener(args[1])

	case "new":
		c.NewListener(args[1:])
		//c.StopListener()
	default:
		c.ui.Send(ui.BAD.Sprintf("Unknown subcommand: %s", args[0]))
	}
}

func (c *CLI) DeleteListeners(name string) {
	endpoint := fmt.Sprintf("ts/rest/listeners/delete/%s", name)
	if err := c.http.DoDelete(endpoint, nil); err != nil {
		c.ui.Send(ui.WARN.Sprintf("Failed Deleting listener: %s", err))
		return
	}
	c.ui.Send(ui.GOOD.Sprintf("Deleted Listener %s", name))
	return

}

func (c *CLI) StopListener(name string) {
	endpoint := fmt.Sprintf("ts/rest/listeners/stop/%s", name)
	if err := c.http.DoPost(endpoint, nil, nil); err != nil {
		c.ui.Send(ui.WARN.Sprintf("Failed Stopping Listener: %s", err))
		return
	}

	c.ui.Send(ui.GOOD.Sprintf("Stopped Listener %s", name))
	return

}

func (c *CLI) ListListners() {
	var r ListListenersResp

	if err := c.http.DoGet("ts/rest/listeners/list", &r); err != nil {
		c.ui.Send(ui.BAD.Sprint("Failed Listing Listeners!"))
		return
	}

	if r.Total == 0 {
		c.ui.Send(ui.INFO.Sprint("No Active Listeners"))
		return
	}

	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "--------\t----\t--------\t----\t------")
	fmt.Fprintln(w, "NAME\tPORT\tPROTOCOL\tHOST\tSTATUS")
	fmt.Fprintln(w, "--------\t----\t--------\t----\t------")

	for _, i := range r.Listeners {
		status := "RUNNING"
		if i.Status == false {
			status = "STOPPED"
		}
		fmt.Fprintf(w, "%s\t%d\t%s\t%s\t%s\n",
			i.Name, i.Port, i.Protocol, i.Host, status)
	}
	w.Flush()
	c.ui.Send(buf.String())
}

func (c *CLI) StartListener(name string) {
	endpoint := fmt.Sprintf("ts/rest/listeners/start/%s", name)
	if err := c.http.DoPost(endpoint, nil, nil); err != nil {
		c.ui.Send(ui.WARN.Sprintf("Failed Deleting listener: %s", err))
		return
	}
	c.ui.Send(ui.GOOD.Sprintf("Started Listener %s", name))
	return
}

func (c *CLI) NewListener(args []string) {
	fs := pflag.NewFlagSet("start", pflag.ContinueOnError)
	port := fs.IntP("port", "p", 443, "")
	proto := fs.StringP("type", "t", "http", "")
	host := fs.StringP("host", "h", "", "")
	cert := fs.BoolP("lets-encrypt", "l", false, "")
	if err := fs.Parse(args); err != nil {
		c.ui.Send(ui.WARN.Sprintf("[!] %v", err))
		return
	}

	if *host == "" {
		c.ui.Send(ui.WARN.Sprint("Must provide host"))
		return
	}

	data := ListenStartReq{
		Port:     *port,
		Protocol: *proto,
		Host:     *host,
		CertType: *cert,
	}

	body, err := json.Marshal(data)
	if err != nil {
		c.ui.Send(ui.BAD.Sprintf("Error Marshaling json: %v", err))
		return
	}
	var StartResp ListenerStartResp
	if err := c.http.DoPost("ts/rest/listeners/new", body, &StartResp); err != nil {
		c.ui.Send(ui.BAD.Sprintf("Error Starting Listener: %v", err))
		return

	}

	c.ui.Send(ui.GOOD.Sprintf("Listener Started: %s", StartResp.Name))
	return

}
