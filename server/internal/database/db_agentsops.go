package database

import (
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func (db *DB) ListAgents() ([]Agent, error) {

	qeuery := `SELECT code_name, username, hostname, external_ip, internal_ip, is_elevated, pid, process_path, windows_version, last_checkin, registration_date FROM agents`

	rows, err := db.conn.Query(qeuery)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	var AgentList []Agent

	for rows.Next() {
		var a Agent
		err := rows.Scan(&a.CodeName, &a.Username, &a.Hostname, &a.Ex_ip, &a.In_ip, &a.IsElevated, &a.Pid, &a.ProcPath, &a.WinVer, &a.LastSeen, &a.RegDate)
		if err != nil {
			return nil, err
		}
		AgentList = append(AgentList, a)
	}

	return AgentList, nil

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

func (db *DB) InsertAgent(guid, codeName, User, Host, InIP, ExIP, ProcPath, WinVer string, Pid int32, IsElev byte) error {
	query := `INSERT INTO agents(guid, code_name,
	  						username, hostname,
							external_ip, internal_ip,
							is_elevated, pid, process_path,
							windows_version, session_key, last_checkin, registration_date) VALUES(
							?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := db.conn.Exec(query, guid, codeName, User, Host, ExIP, InIP, IsElev, Pid, ProcPath, WinVer, "32324234", time.Now().UTC().Unix(), time.Now().UTC().Unix())
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

type AgentInfo struct {
	User         string
	Host         string
	ProcPath     string
	Pid          int32
	WinVer       string
	InternalIP   string
	ExternalIP   string
	IsElevated   bool
	LastCheckin  int64
	RegisterTime int64
}

func (db *DB) ListAgentInfo(name string) (AgentInfo, error) {
	q := `SELECT username, hostname,
			     process_path, pid, windows_version,
				 internal_ip, external_ip,
				 is_elevated, last_checkin, registration_date
		  FROM agents WHERE code_name = ?`

	var a AgentInfo

	err := db.conn.QueryRow(q, name).Scan(&a.User, &a.Host, &a.ProcPath, &a.Pid, &a.WinVer, &a.InternalIP, &a.ExternalIP, &a.IsElevated, &a.LastCheckin, &a.RegisterTime)
	if err != nil {
		fmt.Println(err)

		return AgentInfo{}, err
	}

	return a, nil
}
