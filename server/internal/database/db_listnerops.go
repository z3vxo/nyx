package database

import "fmt"

func (db *DB) InsertListener(port int, id, name, protocol string) error {
	query := `INSERT INTO listeners(port, guid, name, protocol, status) VALUES(?, ?, ?, ?, ?)`

	_, err := db.conn.Exec(query, port, id, name, protocol, "running")
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
}

func (db *DB) GetListeners() ([]ListenersToStart, error) {
	q := `SELECT guid, port, protocol, name FROM listeners WHERE status='running'`

	rows, err := db.conn.Query(q)
	if err != nil {
		fmt.Printf("Failed Query")
		return nil, err
	}
	defer rows.Close()

	var Entrys []ListenersToStart
	for rows.Next() {
		var l ListenersToStart
		err := rows.Scan(&l.Guid, &l.Port, &l.Protocol, &l.Name)
		if err != nil {
			return nil, err
		}

		Entrys = append(Entrys, l)
	}

	return Entrys, nil
}
