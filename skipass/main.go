package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

var (
	Rootfs        = flag.String("rootfs", "", "path to rootfs `dir`")
	WaylandSocket = flag.String("wayland-socket", "", "bind wayland socket to `file`")
	UsrLocalDir   = flag.String("usr-local", "", "bind /usr/local from host to `dir`")
	TmpDir        = flag.String("tmp-dir", "", "bind `dir` in place of /tmp and /run")
	Uid           = flag.Int64("uid", 1000, "change user `id`")
	Gid           = flag.Int64("gid", 1000, "change group `id`")
)

func executable() string {
	filename, _ := os.Executable()
	return filename
}

func init() {
	flag.Parse()

	if os.Getenv("SKIPASS") == "" {
		cmd := exec.Command(executable(), os.Args[1:]...)
		cmd.Env = []string{"SKIPASS=1"}
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		cmd.SysProcAttr = &syscall.SysProcAttr{
			Cloneflags: syscall.CLONE_NEWNS,
		}

		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "skipass: %s\ns", err)
			os.Exit(1)
		}

		os.Exit(0)
	}
}

func main() {
	if err := syscall.Mount("proc", filepath.Join(*Rootfs, "proc"), "proc", 0, ""); err != nil {
		fmt.Fprintf(os.Stderr, "/proc: %s\n", err)
		os.Exit(1)
	}

	if err := syscall.Mount("sysfs", filepath.Join(*Rootfs, "sys"), "sysfs", 0, ""); err != nil {
		fmt.Fprintf(os.Stderr, "/sys: %s\n", err)
		os.Exit(1)
	}

	if err := syscall.Mount("devtmpfs", filepath.Join(*Rootfs, "dev"), "devtmpfs", 0, ""); err != nil {
		fmt.Fprintf(os.Stderr, "/dev: %s\n", err)
		os.Exit(1)
	}

	if err := syscall.Mount("devpts", filepath.Join(*Rootfs, "dev/pts"), "devpts", 0, ""); err != nil {
		fmt.Fprintf(os.Stderr, "/dev/pts: %s\n", err)
		os.Exit(1)
	}

	if *TmpDir != "" {
		base := filepath.Base(*Rootfs)

		tmpDir, err := os.MkdirTemp(filepath.Clean(*TmpDir), base)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", tmpDir, err)
			os.Exit(1)
		}

		runDir, err := os.MkdirTemp(filepath.Clean(*TmpDir), base)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", runDir, err)
			os.Exit(1)
		}

		os.Chmod(tmpDir, 0777)
		os.Chmod(runDir, 0755)

		if err := syscall.Mount(tmpDir, filepath.Join(*Rootfs, "tmp"), "none", syscall.MS_BIND, ""); err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", tmpDir, err)
			os.Exit(1)
		}

		if err := syscall.Mount(runDir, filepath.Join(*Rootfs, "run"), "none", syscall.MS_BIND, ""); err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", runDir, err)
			os.Exit(1)
		}
	} else {
		if err := syscall.Mount("/tmp", filepath.Join(*Rootfs, "tmp"), "none", syscall.MS_BIND, ""); err != nil {
			fmt.Fprintf(os.Stderr, "/tmp: %s\n", err)
			os.Exit(1)
		}

		if err := syscall.Mount("/run", filepath.Join(*Rootfs, "run"), "none", syscall.MS_BIND, ""); err != nil {
			fmt.Fprintf(os.Stderr, "/run: %s\n", err)
			os.Exit(1)
		}
	}

	if *WaylandSocket != "" {
		if err := syscall.Mount("/run/chrome/wayland-0", filepath.Join(*Rootfs, *WaylandSocket), "none", syscall.MS_BIND, ""); err != nil {
			fmt.Fprintf(os.Stderr, "/run/chrome/wayland-0: %s\n", err)
			os.Exit(1)
		}
	}

	if *UsrLocalDir != "" {
		if err := syscall.Mount("/usr/local", filepath.Join(*Rootfs, *UsrLocalDir), "none", syscall.MS_BIND, ""); err != nil {
			fmt.Fprintf(os.Stderr, "/usr/local: %s\n", err)
			os.Exit(1)
		}
	}

	if err := syscall.Chroot(*Rootfs); err != nil {
		fmt.Fprintf(os.Stderr, "chroot: %s\n", err)
		os.Exit(1)
	}

	os.Chdir("/")

	args := flag.Args()
	if len(args) == 0 {
		args = []string{"/bin/sh"}
	}

	if err := syscall.Setgid(int(*Gid)); err != nil {
		fmt.Fprintf(os.Stderr, "gid: %s\n", err)
		os.Exit(1)
	}

	if err := syscall.Setuid(int(*Uid)); err != nil {
		fmt.Fprintf(os.Stderr, "uid: %s\n", err)
		os.Exit(1)
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", args[0], err)
		os.Exit(1)
	}
}
