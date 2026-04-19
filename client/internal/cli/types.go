package cli

type Agent struct {
	CodeName   string `json:"code_name"`
	Username   string `json:"username"`
	Hostname   string `json:"hostname"`
	Ex_ip      string `json:"ex_ip"`
	In_ip      string `json:"in_ip"`
	IsElevated bool   `json:"is_elevated"`
	Pid        int    `json:"pid"`
	ProcPath   string `json:"proc_path"`
	WinVer     string `json:"winver"`
	LastSeen   int64  `json:"last_checkin"`
}

type Agents struct {
	Total int     `json:"total"`
	Agent []Agent `json:"agents"`
}

type ResolveResp struct {
	Guid string `json:"guid"`
}

type ListenerEntry struct {
	Port   int
	Name   string
	Status string
}

type ListListenersResp struct {
	Total     int             `json:"total"`
	Listeners []ListenerEntry `json:"listeners"`
}
