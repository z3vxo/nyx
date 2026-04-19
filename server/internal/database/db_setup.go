package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	conn *sql.DB
}

func GetDbPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/.kronos/database/kronos_db.sql", home), nil
}

func NewDB() (*DB, error) {
	dbPath, err := GetDbPath()
	if err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	d := &DB{conn: db}
	err = InitDB(d)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func InitDB(db *DB) error {

	pragmans := []string{
		"PRAGMA journal_mode=WAL",
		"PRAGMA synchronous=NORMAL",
		"PRAGMA foreign_keys=ON",
	}

	for _, p := range pragmans {
		if _, err := db.conn.Exec(p); err != nil {
			fmt.Println("pragmas")
			return err
		}
	}

	return SetupDB(db)
}

func SetupDB(db *DB) error {

	agents_querys := `CREATE TABLE IF NOT EXISTS agents (
		guid 			TEXT NOT NULL,
		code_name 		TEXT NOT NULL,
		username 		TEXT NOT NULL,
		hostname 		TEXT NOT NULL,
		external_ip	 	TEXT NOT NULL,
		internal_ip 	TEXT NOT NULL,
		is_elevated 	BOOLEAN NOT NULL,
		pid 			INTEGER NOT NULL,
		process_path 	TEXT NOT NULL,
		windows_version TEXT NOT NULL,
		session_key    	BLOB NOT NULL,
		last_checkin    INTEGER NOT NULL);`

	_, err := db.conn.Exec(agents_querys)
	if err != nil {
		fmt.Println(err)
		fmt.Println("agents")
		return err
	}

	commands_query := `CREATE TABLE IF NOT EXISTS commands (
		guid TEXT NOT NULL,
		command_type INTEGER NOT NULL,
		task_id      INTEGER NOT NULL,
		param_1      TEXT NOT NULL,
		param_2      TEXT NOT NULL,
		executed     BOOLEAN NOT NULL,
		tasked_at    INTEGER NOT NULL);`

	_, err = db.conn.Exec(commands_query)
	if err != nil {

		return err
	}

	listeners_query := `CREATE TABLE IF NOT EXISTS listeners (
		id INTEGER PRIMARY KEY,
		guid TEXT NOT NULL,
		port INTEGER NOT NULL,
		name TEXT NOT NULL,
		status TEXT NOT NULL);
		`
	_, err = db.conn.Exec(listeners_query)
	if err != nil {
		return err
	}

	return nil
}
