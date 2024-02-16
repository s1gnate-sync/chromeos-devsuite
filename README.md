![Screenshot 2024-02-16 17 41 21](https://github.com/s1gnate-sync/chromeos-devsuite/assets/139636216/a30dda10-5530-4485-8a81-589e94569677)

```
root      4568  0.0  0.0 1226900    0 ?        Ssl  17:07   0:00 /usr/local/SUITE/bin/devcoo once { /usr/local/SUITE/bin
root      4714  0.0  0.0 1229128 3032 ?        Sl   17:07   0:01  \_ /usr/local/SUITE/bin/sucrosh -addr 127.0.0.22:22 -k
chronos   8438  0.0  0.0   3848  3048 pts/5    Ss+  17:20   0:00  |   \_ /bin/bash -l
root      4715  0.0  0.0 1226688    0 ?        Sl   17:07   0:00  \_ /usr/local/SUITE/bin/skipass -rootfs /usr/local/ski
root      4735  0.0  0.0 1226688   68 ?        Sl   17:07   0:00  |   \_ /usr/local/SUITE/bin/skipass -rootfs /usr/local
root      4744  0.1  0.0 1229128 3216 ?        Sl   17:07   0:02  |       \_ /opt/SUITE/bin/sucrosh -addr 127.0.0.222:22
chronos   9105  0.0  0.0   2940  2468 pts/1    Ss+  17:23   0:00  |           \_ /bin/bash -l
root      4717  0.0  0.0 1226688   24 ?        Sl   17:07   0:00  \_ /usr/local/SUITE/bin/skipass -rootfs /usr/local/ski
root      4736  0.0  0.0 1226944   44 ?        Sl   17:07   0:00      \_ /usr/local/SUITE/bin/skipass -rootfs /usr/local
root      4745  0.0  0.0 1229832 3528 ?        Sl   17:07   0:01          \_ /opt/SUITE/bin/sucrosh -addr 127.0.0.220:22
chronos  12304  0.0  0.0   4052  3240 pts/0    Ss+  17:39   0:00              \_ /bin/bash -l
```

# chromeos-devsuite
Tools to improve devmode experience: privileged shell, upstart integration, hosts file management, tools for chroot environments with namespaces

# rootfs

for debian bookworm and alpine stable I simply ripped images from docker and created chronos user

# installation

to compile simply write make and make sure go is installed

should work out of the box if placed into hardcoded path `/usr/local/suite/`
