#!/bin/bash

# This is a simple script to run our tests and accumulate the
# resulting coverage report into a single coverage.out file

fail=0

echo "mode: set" > coverage.out

go test -v . -covermode=count -coverprofile=magic_packet.out || fail=1
cat magic_packet.out | tail -n +2 >> coverage.out
rm magic_packet.out

go test -v ./cmd/wol -covermode=count -coverprofile=wol.out || fail=1
cat wol.out | tail -n +2 >> coverage.out
rm wol.out
