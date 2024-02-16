#!/bin/sh
'': 
	CGO_ENABLED=0 exec go build -buildvcs=false -ldflags="-s -w" -trimpath
.PHONY: ''

