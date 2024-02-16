#!/bin/sh

# set hosts and run priveledged shell
exec /usr/local/suite/bin/devcoo \
	once { /usr/local/suite/bin/sethosts } \
		 { /usr/local/suite/conf/sucrosh-host.sh } \
		 { /usr/local/suite/conf/sucrosh-skipass-bookworm.sh } \
		 { /usr/local/suite/conf/sucrosh-skipass-alpine.sh }
