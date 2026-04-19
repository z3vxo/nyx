#!/bin/bash


token=$(curl http://localhost:50050/ts/rest/login -s -d '{"username":"kronos", "password":"kronospwd"}' | jq -r '.token')
echo "$token"

#curl -s http://localhost:3000/rest/agents -s -H "Authorization: Bearer $token" | jq

# type TaskEntry struct {
# 	Cmd_type int    `json:"type"`
# 	Guid     string `json:"guid"`
# 	Param1   string `json:"param1"`
# 	Param2   string `json:"param2"`
# }
case "$1" in
    resolve)
        curl -s "http://localhost:50050/ts/rest/agents/resolve/$2"  -H "Authorization: Bearer $token" | jq
        ;;
    list)
        curl -s http://localhost:50050/ts/rest/agents/list -s -H "Authorization: Bearer $token" | jq
        ;;
    insert)
        curl -s http://localhost:50050/ts/rest/commands/new  -d '{"guid":"1122", "type":1, "task_id":"1111-11111", "param_1":"test"}' -s -H "Authorization: Bearer $token" | jq
        ;;
    delete_cmd)
        curl -s http://localhost:50050/ts/rest/commands/delete  -d '{"task_id":""}' -s -H "Authorization: Bearer $token" | jq
        ;;
    list_cmd)
        curl -s "http://localhost:50050/ts/rest/tasks/list/$2" -s -H "Authorization: Bearer $token" | jq
        ;;
    list_start)
        curl -s "http://localhost:50050/ts/rest/listeners/start" -X POST -d '{"port":8081}' -H "Authorization: Bearer $token" | jq
        ;;
    list_stop)
        curl -s "http://localhost:50050/ts/rest/listeners/stop/$2" -X POST -H "Authorization: Bearer $token" | jq
        ;;
    sse)
        curl -N http://localhost:50050/ts/events -H "Authorization: Bearer $token"
        ;;
    *)

        echo "need option"
esac
