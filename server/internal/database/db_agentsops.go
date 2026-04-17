package database

import (
	"encoding/json"

	_ "github.com/mattn/go-sqlite3"
)

func (db *DB) ListAgents() ([]byte, error) {

	qeuery := `SELECT code_name, username, hostname, external_ip, internal_ip, is_elevated, pid, process_path, windows_version, last_checkin FROM agents`

	rows, err := db.conn.Query(qeuery)
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

func (db *DB) ResolveCodename(name string) (string, error) {
	query := `SELECT guid FROM agents WHERE code_name = ?`
	var guid string
	err := db.conn.QueryRow(query, name).Scan(&guid)

	if err != nil {
		return "", err
	}

	return guid, nil
}
