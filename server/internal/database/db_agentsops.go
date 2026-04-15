package database

import (
	"encoding/json"

	_ "github.com/mattn/go-sqlite3"
)

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
	LastSeen   string `json:"last_checkin"`
}

type Agents struct {
	Total int     `json:"total"`
	Agent []Agent `json:"agents"`
}

func Db_ListAgents() ([]byte, error) {

	qeuery := `SELECT code_name, username, hostname, external_ip, internal_ip, is_elevated, pid, process_path, windows_version, last_checkin FROM agents`

	rows, err := db.Query(qeuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var AgentList []Agent

	for rows.Next() {
		var a Agent
		err := rows.Scan(&a.CodeName, &a.Username, &a.Hostname, &a.Ex_ip, &a.In_ip, &a.IsElevated, &a.Pid, &a.ProcPath, &a.WinVer, &a.LastSeen)
		if err != nil {
			return nil, err
		}
		AgentList = append(AgentList, a)
	}

	payload := Agents{
		Total: len(AgentList),
		Agent: AgentList,
	}

	res, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return res, nil

}

func ResolveCodename(name string) (string, error) {
	query := `SELECT guid FROM agents WHERE code_name = ?`
	var guid string
	err := db.QueryRow(query, name).Scan(&guid)

	if err != nil {
		return "", err
	}

	return guid, nil
}
