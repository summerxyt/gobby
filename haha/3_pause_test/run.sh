#!/bin/sh
rm -f dumplog_*
go run node1/node.go &
go run node2/node.go &
go run node3/node.go &
