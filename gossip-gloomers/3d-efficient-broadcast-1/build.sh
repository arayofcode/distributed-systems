#!/bin/bash

MAELSTROM_PATH=../maelstrom
cwd=$(pwd)
go build
cd $MAELSTROM_PATH
./maelstrom test -w broadcast --bin $cwd/maelstrom-broadcast --node-count 25 --time-limit 20 --rate 100 --latency 100  --nemesis partition
cd $cwd
rm $cwd/maelstrom-broadcast