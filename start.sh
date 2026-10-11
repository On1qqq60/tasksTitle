#!/bin/bash
set -euo pipefail

docker-entrypoint.sh postgres &
postgres_pid=$!

for _ in $(seq 1 90); do
	if psql -U "${POSTGRES_USER}" -d "${POSTGRES_DB}" -c "SELECT 1" >/dev/null 2>&1; then
		break
	fi
	if ! kill -0 "$postgres_pid" 2>/dev/null; then
		echo "Postgres завершился" >&2
		exit 1
	fi
	sleep 1
done

if ! psql -U "${POSTGRES_USER}" -d "${POSTGRES_DB}" -c "SELECT 1" >/dev/null 2>&1; then
	echo "База ${POSTGRES_DB} недоступна" >&2
	exit 1
fi

/usr/local/bin/notes &
app_pid=$!

wait -n "$postgres_pid" "$app_pid"
exit $?
