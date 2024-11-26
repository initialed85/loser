#!/bin/bash

set -e

docker build -t loser-iproute2 -f ./Dockerfile.iproute2 .

if [[ "${INTERNAL}" != "1" ]]; then
    docker run --rm -it --name loser-iproute2 -v "$(pwd):/srv" -p 6942:6942/tcp -p 6943:6943/tcp -p 6943:6943/udp --entrypoint bash loser-iproute2
else
    docker run --rm -it --name loser-iproute2-internal -v "$(pwd):/srv" --entrypoint bash loser-iproute2
fi
