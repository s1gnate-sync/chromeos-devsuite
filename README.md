![Screenshot 2024-02-16 17 41 21](https://github.com/s1gnate-sync/chromeos-devsuite/assets/139636216/a30dda10-5530-4485-8a81-589e94569677)
# chromeos-devsuite
Tools to improve devmode experience: privileged shell, upstart integration, hosts file management, tools for chroot environments with namespaces

# rootfs

for debian bookworm and alpine stable I simply ripped images from docker and created chronos user

# installation

to compile simply write make and make sure go is installed

should work out of the box if placed into hardcoded path `/usr/local/suite/`
