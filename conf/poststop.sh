#!/bin/sh

/bin/umount /etc/sudoers.d 2>/dev/null \
	|| /bin/true

rm -fr /usr/local/tmp
