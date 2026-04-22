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
	RegDate    int64  `json:"reg_date"`
}

type AgentInfoResp struct {
	User         string `json:"username"`
	Host         string `json:"hostname"`
	ProcPath     string `json:"proc_path`
	Pid          int32  `json:"pid"`
	WinVer       string `json:"win_version`
	InternalIP   string `json:"internal_ip"`
	ExternalIP   string `json:"external_ip"`
	IsElevated   bool   `json:"is_elev`
	LastCheckin  int32  `json:"last_checkin"`
	RegisterTime int32  `json:"reg_date`
}

type Agents struct {
	Total int     `json:"total"`
	Agent []Agent `json:"agents"`
}

type ResolveResp struct {
	Guid string `json:"guid"`
}

// ----- Listener List response data -----
type ListenerEntry struct {
	Port     int
	Name     string
	Protocol string
	Status   bool
	Host     string
}

type ListListenersResp struct {
	Total     int             `json:"total"`
	Listeners []ListenerEntry `json:"listeners"`
}

//----- listener Start request data -----

type ListenStartReq struct {
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
	Host     string `json:"host"`
	CertType bool   `json:"letsencrypt"` // 0 = self signed, 1 = lets encrypt
}

// Listener Start Response data
type ListenerStartResp struct {
	Name string `json:"listener_name"`
}

type Generic200 struct {
	Status string `json:"status"`
}
