#!/bin/bash

MAELSTROM_PATH=../maelstrom
cwd=$(pwd)
go build -o bin
cd $MAELSTROM_PATH
./maelstrom test -w broadcast --bin $cwd/bin --node-count 5 --time-limit 20 --rate 10
cd $cwd