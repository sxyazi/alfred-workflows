#!/usr/bin/env bash
cd "$(dirname "$0")"
go build -o dist/alfred-workflows .
./dist/alfred-workflows "$1" "$2"
