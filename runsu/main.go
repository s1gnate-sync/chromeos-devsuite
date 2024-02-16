package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"syscall"
)

var (
	Sh 		= flag.String("sh", "/bin/sh", "path to sh command")
	Enable  = flag.Bool("su-enable", false, "allow regular users run root commands")
	Disable = flag.Bool("su-disable", false, "clear command from extra capabilities")
	Uid     = flag.Uint("uid", 0, "change user `id`")
	Gid     = flag.Uint("gid", 0, "change group `id`")
)

func init() {
	flag.Parse()
	executable, _ := os.Executable()
	if executable == "" {
		executable = os.Args[0]
	}

	if *Enable {
		syscall.Chown(executable, 0, 0)
		syscall.Chmod(executable, syscall.S_IRWXU | syscall.S_IXOTH | syscall.S_IXGRP | syscall.S_ISUID)
	} else if *Disable {
		syscall.Chown(executable, 65534, 65534)
	}

	if *Enable || *Disable {
		os.Exit(0)
	}
}

const hint = "capabilities must be enabled for this command to work"

func main() {
	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "too few arguments\n")
		os.Exit(1)
	}

	if err := syscall.Setuid(0); err != nil {
		fmt.Fprintf(os.Stderr, "setuid: %s: %s (%s)\n", os.Args[0], hint, err)
		os.Exit(1)
		
	}

	if err := syscall.Setgid(0); err != nil {
		fmt.Fprintf(os.Stderr, "setgid: %s: %s (%s)\n", os.Args[0], hint, err)
		os.Exit(1)		
	}

	pid, err := syscall.ForkExec(*Sh, []string{"sh", "-l", "-c", strings.Join(args, " ")}, &syscall.ProcAttr{
		Files: []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()},
		Sys: &syscall.SysProcAttr{
				Credential: &syscall.Credential{
					Uid : uint32(*Uid),
					Gid: uint32(*Gid),
				},
			},
	})

	if err != nil {
		fmt.Fprintf(os.Stdin, "exec: %s: %s\n", args[0], err)
		os.Exit(1)
	}
	
	var status syscall.WaitStatus
	syscall.Wait4(pid, &status, 0, nil)
	if !status.Exited() {
		fmt.Fprintf(os.Stderr, "%s: status %s\n", args, status.ExitStatus())
	}
	
	os.Exit(status.ExitStatus())	
}
