#!/bin/sh

set -e
set -x

docker build -t ropoly-cmd .
docker run --rm -it -v $PWD:/out ropoly-cmd bash -c "cp ropoly-cmd /out/ropoly-cmd"