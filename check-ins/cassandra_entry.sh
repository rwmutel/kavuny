#!/usr/bin/bash

until cat /opt/app/create_check_ins.cql | cqlsh; do
echo "cqlsh: Cassandra is unavailable - retry later"
sleep 5
done &

echo "cqlsh: Cassandra is populated"
bash docker-entrypoint.sh "$@"
