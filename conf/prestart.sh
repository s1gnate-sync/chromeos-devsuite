#!/bin/sh

mkdir /usr/local/tmp
chmod 0777 /usr/local/tmp

# replace sudoers
/bin/mountpoint -q /etc/sudoers.d \
	|| /bin/mount --bind /usr/local/suite/conf/sudoers.d /etc/sudoers.d

for n in 22 222 220; do
	/bin/ip addr add 127.0.0.$n dev lo 2>/dev/null \
		|| /bin/true
done

/usr/local/suite/bin/runsu -su-enable || /bin/true
