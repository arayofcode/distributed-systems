#!/bin/bash

MAELSTROM_PATH=../maelstrom
cwd=$(pwd)
go build -o bin
cd $MAELSTROM_PATH
./maelstrom test -w broadcast --bin $cwd/ --node-count 25 --time-limit 20 --rate 100 --latency 100
cd $cwd