package database

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

type Agents struct {
	Total int     `json:"total"`
	Agent []Agent `json:"agents"`
}

type Task struct {
	CmdCode  int    `json:"cmd_code"`
	Param1   string `json:"param_1"`
	Param2   string `json:"param_2"`
	TaskedAt int    `json:"tasked_at"`
}

type TaskEntrys struct {
	Total int    `json:"total"`
	Tasks []Task `json:"tasks"`
}
