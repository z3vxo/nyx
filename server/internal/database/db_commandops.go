package database

import (
	"fmt"
	"time"
)

func InsertCommand(cmdType int, guid, param1, param2 string) error {

	query := `INSERT INTO commands(guid, command_type, param_1, param_2, executed, tasked_at) VALUES(?, ?, ?, ?, ?, ?)`

	_, err := db.Exec(query, guid, cmdType, param1, param2, 0, time.Now().Unix())
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil

}
