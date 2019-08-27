#!/bin/bash

if [ $# -lt 1 ]; then
    echo "Usage ./build.sh [TAG]"
    exit 1
fi

docker build . -t wxproxy:latest
docker tag wxproxy:latest yuezzzzzzzzzzz/wx-proxy:latest
docker tag wxproxy:latest yuezzzzzzzzzzz/wx-proxy:$1

docker push yuezzzzzzzzzzz/wx-proxy:latest
docker push yuezzzzzzzzzzz/wx-proxy:$1
