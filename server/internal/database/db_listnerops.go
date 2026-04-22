package database

import (
	"fmt"
)

func (db *DB) InsertListener(port int, id, name, protocol, host string, certType, status bool) error {
	query := `INSERT INTO listeners(port, guid, name, protocol, host, certType, status) VALUES(?, ?, ?, ?, ?, ?, ?)`

	_, err := db.conn.Exec(query, port, id, name, protocol, host, certType, status)
	if err != nil {
		fmt.Printf("Failed Insert: %v", err)
		return err
	}

	return nil
}

func (db *DB) DeleteListener(id string) error {
	query := `DELETE FROM listeners WHERE guid = ?`
	_, err := db.conn.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

type ListenersToStart struct {
	Guid     string
	Port     int
	Name     string
	Protocol string
	Host     string
	CertType bool
	Status   bool
}

func (db *DB) GetListeners() ([]ListenersToStart, error) {
	q := `SELECT guid, port, protocol, name, host, certType, status FROM listeners`

	rows, err := db.conn.Query(q)
	if err != nil {
		fmt.Printf("Failed Query")
		return nil, err
	}
	defer rows.Close()

	var Entrys []ListenersToStart
	for rows.Next() {
		var l ListenersToStart
		err := rows.Scan(&l.Guid, &l.Port, &l.Protocol, &l.Name, &l.Host, &l.CertType, &l.Status)
		if err != nil {
			return nil, err
		}

		Entrys = append(Entrys, l)
	}

	return Entrys, nil
}

func (db *DB) UpdateListenerStatus(id string, status bool) error {
	q := `UPDATE listeners SET status = ? WHERE guid = ?`
	_, err := db.conn.Exec(q, status, id)
	if err != nil {
		return err
	}
	return nil
}
