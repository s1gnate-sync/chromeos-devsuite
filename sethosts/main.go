package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	File string
)

func executable() string {
	filename, _ := os.Executable()
	return filename
}

func confFile(name string) string {
	return filepath.Clean(filepath.Join(filepath.Dir(executable()), "..", "conf", name))
}

func main() {
	flag.StringVar(&File, "file", confFile("hosts"), "`path` to hosts file")
	flag.Parse()

	data, err := os.ReadFile(File)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read hosts: %s\n", err)
		os.Exit(1)
	}

	if err := exec.Command("/sbin/stop", "crosdns").Run(); err != nil {
		fmt.Fprintf(os.Stderr, "stop crosdns: %s\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile("/run/crosdns/hosts", data, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "write hosts: %s\n", err)
	}

	if err := exec.Command("/sbin/start", "crosdns").Run(); err != nil {
		fmt.Fprintf(os.Stderr, "start crosdns: %s\n", err)
		os.Exit(1)
	}
}
