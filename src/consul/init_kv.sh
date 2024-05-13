#!/bin/sh

echo "bootstraping values"
echo "waiting until port is available..."
while ! nc -z localhost 8500; do   
  sleep 1
done

echo "converting values to base64 and executing consul kv command"
consul kv import "$(cat /opt/consul_kv.json | jq 'walk(if type == "object" then with_entries(if .key == "value" then .value |= @base64 else . end) else . end)')"
