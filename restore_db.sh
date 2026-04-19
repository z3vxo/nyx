#!/usr/bin/env bash
set -e

DB="$HOME/.kronos/database/kronos_db.sql"

sqlite3 "$DB" <<'EOF'

-- agents
INSERT INTO agents (guid, code_name, username, hostname, external_ip, internal_ip, is_elevated, pid, process_path, windows_version, session_key, last_checkin) VALUES
  ('1a2b3c4d-1234-5678-abcd-ef0123456789', 'ted_kaz', 'DESKTOP-ABC\john', 'DESKTOP-ABC', '1.2.3.4', '192.168.1.10', 1, 4812, 'C:\Users\john\AppData\Local\svchost.exe', 'Windows 10 Pro 22H2', X'DEADBEEF', 1776402339),
  ('2b3c4d5e-2345-6789-bcde-f01234567890', 'iron_fox', 'CORP-PC\administrator', 'CORP-PC', '5.6.7.8', '10.0.0.5', 1, 1337, 'C:\Windows\System32\svchost.exe', 'Windows 11 Pro 23H2', X'CAFEBABE', 1776402372),
  ('3c4d5e6f-3456-789a-cdef-012345678901', 'dark_owl', 'WORKSTATION-7\msmith', 'WORKSTATION-7', '9.10.11.12', '172.16.0.22', 0, 9204, 'C:\Users\msmith\Downloads\update.exe', 'Windows 10 Home 21H2', X'BEEFDEAD', 1776402372);

-- commands
INSERT INTO commands (guid, command_type, task_id, param_1, param_2, executed, tasked_at) VALUES
  ('1a2b3c4d-1234-5678-abcd-ef0123456789', 1, 'task-0001', 'whoami', '', 0, 1776402372),
  ('1a2b3c4d-1234-5678-abcd-ef0123456789', 1, 'task-0002', 'ipconfig /all', '', 1, 1776402072),
  ('2b3c4d5e-2345-6789-bcde-f01234567890', 2, 'task-0003', 'C:\Users\administrator\secret.txt', '', 0, 1776402372),
  ('3c4d5e6f-3456-789a-cdef-012345678901', 1, 'task-0004', 'net user', '', 0, 1776402372),
  ('test-guid-1234-abcd', 1, 'task-001', 'ls', '-la', 0, 1776445202),
  ('test-guid-1234-abcd', 2, 'task-002', 'whoami', '', 0, 1776445202);

-- listeners (new schema includes name column)
INSERT INTO listeners (guid, port, name, protocol, status) VALUES
  ('22ec4e97-40f8-4183-9247-c662504518c8', 8080, 'restored-listener', 'http', 'running');

EOF

echo "DB restored: $DB"
