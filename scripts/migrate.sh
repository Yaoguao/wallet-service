#!/bin/sh

if [ "$1" = "down" ]; then
  shift
  migrate -path=/app/migrations -database=$DSN_POSTGRES down "${1:-1}"
else
  migrate -path=/app/migrations -database=$DSN_POSTGRES up
fi