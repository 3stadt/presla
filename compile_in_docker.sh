#!/usr/bin/env bash

## call this script with argument 'RELEASE_VERSION=vx.x.x' for compiling another version than "vlatest"

docker build -t preslago:1.10 ./docker/

UG=$(id -g "$USER")

docker run --rm -it -v "$PWD"/docker:/root/.config -v "$PWD":/go/src/github.com/3stadt/presla -w /go/src/github.com/3stadt/presla preslago:1.10 bash -c "make bindata-debug && make all $1 || chown -R $UID:$UG ./"