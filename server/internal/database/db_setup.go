package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func GetDbPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/.kronos/database/kronos_db.sql", home), nil
}

func InitDB() error {
	var err error
	dbPath, err := GetDbPath()
	if err != nil {
		fmt.Println("path")
		return err
	}
	fmt.Println(dbPath)

	DB, err := sql.Open("sqlite3", dbPath)

	if err := DB.Ping(); err != nil {
		fmt.Println("ping")

		return err
	}

	pragmans := []string{
		"PRAGMA journal_mode=WAL",
		"PRAGMA synchronous=NORMAL",
		"PRAGMA foreign_key=ON",
	}

	for _, p := range pragmans {
		if _, err := DB.Exec(p); err != nil {
			fmt.Println("pragmas")
			return err
		}
	}
	db = DB

	return SetupDB()
}

func SetupDB() error {

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

	_, err := db.Exec(agents_querys)
	if err != nil {
		fmt.Println(err)
		fmt.Println("agents")
		return err
	}

	commands_query := `CREATE TABLE IF NOT EXISTS commands (
		guid TEXT NOT NULL,
		command_type INTEGER NOT NULL,
		task_id      TEXT NOT NULL,
		param_1      TEXT NOT NULL,
		param_2      TEXT NOT NULL,
		executed     BOOLEAN NOT NULL,
		tasked_at    INTEGER NOT NULL);`

	_, err = db.Exec(commands_query)
	if err != nil {
		fmt.Println(err)
		fmt.Println("commands")

		return err
	}

	return nil
}
