#!/bin/sh

/opt/init_kv.sh &

consul agent -server -ui -node=server-1 -bootstrap-expect=1 -client=0.0.0.0 --data-dir /consul/data
