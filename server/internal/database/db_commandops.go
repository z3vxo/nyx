package database

import (
	"fmt"
	"time"
)

func (db *DB) InsertCommand(cmdType int, guid, taskid, param1, param2 string) error {

	query := `INSERT INTO commands(guid, command_type, task_id, param_1, param_2,executed, tasked_at) VALUES(?, ?, ?, ?, ?, ?, ?)`

	_, err := db.conn.Exec(query, guid, cmdType, taskid, param1, param2, 0, time.Now().Unix())
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil

}

func (db *DB) DeleteTask(id int) error {
	query := `DELETE FROM commands WHERE task_id = ?`

	_, err := db.conn.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
