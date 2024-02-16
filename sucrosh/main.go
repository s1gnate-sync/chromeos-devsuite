package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sucrosh/ssh"
	"syscall"
)

var (
	Address             = flag.String("addr", "127.0.0.1:2222", "listen on specified `address:port`")
	HostKey      string = ""
	User         string = "chronos"
	Uid                 = flag.Int64("uid", 1000, "change user `id`")
	Gid                 = flag.Int64("gid", 1000, "change group `id`")
	Dir          string = ""
	UnshareMount        = flag.Bool("unshare-mount", false, "unshare mount namespace")
	Args                = []string{}
	Env                 = flag.String("env", "", "comma separated list of environment variables `NAME=VAL,NAME2=VAL2`")
)

func executable() string {
	filename, _ := os.Executable()
	return filename
}

func confFile(name string) string {
	return filepath.Clean(filepath.Join(filepath.Dir(executable()), "..", "conf", name))
}

func init() {
	whoami, _ := exec.Command("whoami").Output()
	user := strings.TrimSpace(string(whoami))
	if user == "" {
		user = "user"
	}

	dir := confFile("../..")
	key := confFile("sucrosh_host_key")

	flag.StringVar(&User, "user", user, "`name` of the user")
	flag.StringVar(&Dir, "dir", dir, "change working directory to `dir`")
	flag.StringVar(&HostKey, "key", key, "host key")

	flag.Parse()
	Args = flag.Args()
	if len(Args) == 0 {
		Args = []string{"/bin/bash", "-l"}
	}

	if os.Getenv("UNSHARED") == "" && *UnshareMount {
		err := (&exec.Cmd{
			Path: executable(),
			Args: os.Args,
			SysProcAttr: &syscall.SysProcAttr{
				Cloneflags: syscall.CLONE_NEWNS,
			},
			Env:    []string{"UNSHARED=1"},
			Stdout: os.Stdout,
			Stdin:  os.Stdin,
			Stderr: os.Stderr,
		}).Run()

		if err != nil {
			fmt.Fprintf(os.Stderr, "unshare mount: %s\n", err)
			os.Exit(1)
		}

		os.Exit(0)
	}
}

func execCommand(name string, args []string) *exec.Cmd {
	if len(args) == 0 {
		args = Args
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = Dir

	env := []string{}
	for _, item := range strings.Split(*Env, ",") {
		env = append(env, item)
	}

	cmd.Env = env

	if os.Getuid() == 0 {
		if name == "root" {
			cmd.SysProcAttr = &syscall.SysProcAttr{
				Credential: &syscall.Credential{
					Uid: 0,
					Gid: 0,
				},
			}
		} else {
			if *Uid == 0 || *Gid == 0 {
				fmt.Fprintf(
					os.Stderr,
					"%s: can't provide command for normal user with uid %d gid %d\n",
					name, *Uid, *Gid,
				)
				return nil
			}

			cmd.SysProcAttr = &syscall.SysProcAttr{
				Credential: &syscall.Credential{
					Uid: uint32(*Uid),
					Gid: uint32(*Gid),
				},
			}
		}
	}

	return cmd
}

func main() {
	fmt.Fprintf(os.Stderr, "running as %s, listening on %s \n", User, *Address)

	err := ssh.Run(*Address, HostKey, map[string]ssh.CmdProvider{
		User: execCommand,
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "run: %s\n", err)
		return
	}
}
