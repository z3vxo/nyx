package database

func (db *DB) InsertListener(port int, id string) error {
	query := `INSERT INTO listeners(port, guid, status) VALUES(?, ?, ?)`

	_, err := db.conn.Exec(query, port, id, "running")
	if err != nil {
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
	Guid string
	Port int
}

func (db *DB) GetListeners() ([]ListenersToStart, error) {
	q := `SELECT guid, port FROM listeners WHERE status='running'`

	rows, err := db.conn.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var Entrys []ListenersToStart
	for rows.Next() {
		var l ListenersToStart
		err := rows.Scan(&l.Guid, &l.Port)
		if err != nil {
			return nil, err
		}

		Entrys = append(Entrys, l)
	}

	return Entrys, nil
}
