#!/bin/sh

set -eux
#	-tmp-dir '/usr/local/tmp' \

# container shell
exec /usr/local/suite/bin/skipass \
	-rootfs '/usr/local/skipass-alpine' \
	-uid '0' \
	-gid '0' \
	-tmp-dir '/usr/local/tmp' \
	-usr-local '/opt' \
	-- /opt/suite/bin/sucrosh \
		-addr '127.0.0.222:22' \
		-key '/opt/suite/conf/sucrosh_skipass_alpine_key' \
		-uid '1000' \
		-gid '1000' \
		-user 'chronos' \
		-env 'HOME=/home/chronos/user,PATH=/bin:/usr/bin:/sbin:/usr/sbin:/usr/local/bin:/opt/suite/bin' \
		-dir '/' 
