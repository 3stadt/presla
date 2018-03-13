#!/usr/bin/env bash

SLIDES="$1"

docker run --rm -t --net=host --shm-size 2G -v "$(pwd)":/slides astefanutti/decktape http://localhost:8080/"$SLIDES" "$SLIDES".pdf