#!/usr/bin/env bash

UG=$(id -g "$USER")

docker run --rm -it -p 8080:8080 -v "$PWD":/go/src/github.com/3stadt/presla -w /go/src/github.com/3stadt/presla golang:1.10 bash -c "make run || chown -R $UID:$UG ./"