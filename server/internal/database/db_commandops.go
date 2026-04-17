package database

import (
	"fmt"
	"time"
)

func (db *DB) InsertCommand(cmdType, taskid int, guid, param1, param2 string) error {

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

func (db *DB) ListTasks(guid string) ([]Task, error) {
	query := `SELECT command_type, param_1, param_2, tasked_at FROM commands WHERE guid = ? AND executed = 0`

	rows, err := db.conn.Query(query, guid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.CmdCode, &t.Param1, &t.Param2, &t.TaskedAt)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	return tasks, nil

}
