#!/bin/sh

# chrome os shell
exec /usr/local/suite/bin/sucrosh \
	-addr '127.0.0.22:22' \
	-key '/usr/local/suite/conf/sucrosh_host_key' \
	-user 'chronos' \
	-uid '1000' \
	-gid '1000' \
	-dir '/' \
	-env 'TMPDIR=/usr/local/tmp,HOME=/home/chronos/user,PATH=/bin:/usr/bin:/sbin:/usr/sbin:/usr/local/bin' \
	-- '/bin/bash' '-l'
