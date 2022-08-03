#!/bin/sh

go mod download && go mod verify
go build -v -o /usr/local/bin/app ./...
# start cron
/usr/sbin/crond -f -l 8
