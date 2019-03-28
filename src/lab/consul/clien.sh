#!/bin/bash

docker run -d --name=node4  --restart=always \
    -e 'CONSUL_LOCAL_CONFIG={"leave_on_terminate": true}' \
    -p 11300:8300 \
    -p 11301:8301 \
    -p 11301:8301/udp \
    -p 11302:8302/udp \
    -p 11302:8302 \
    -p 11400:8400 \
    -p 11500:8500 \
    -p 11600:8600 \
    -h node4 \
    consul agent -bind=172.17.0.5 -retry-join=192.168.83.129 \
    -node-id=$(uuidgen | awk '{print tolower($0)}') \
    -node=node4 -client 0.0.0.0 -ui

