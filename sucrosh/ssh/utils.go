package ssh

import (
	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
	"os"
	"syscall"
	"unsafe"
)

func setWinsize(f *os.File, w, h int) {
	syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), uintptr(syscall.TIOCSWINSZ),
		uintptr(unsafe.Pointer(&struct{ h, w, x, y uint16 }{uint16(h), uint16(w), 0, 0})))
}

func keyboardInteractiveHandler(ctx ssh.Context, challenger gossh.KeyboardInteractiveChallenge) bool {
	return true
}

func ptyCallback(ctx ssh.Context, pty ssh.Pty) bool {
	return true
}

func publicKeyHandler(ctx ssh.Context, key ssh.PublicKey) bool {
	return false
}

func passwordHandler(ctx ssh.Context, password string) bool {
	return false
}

func localPortForwardingCallback(ctx ssh.Context, destinationHost string, destinationPort uint32) bool {
	return false
}

func reversePortForwardingCallback(ctx ssh.Context, bindHost string, bindPort uint32) bool {
	return false
}
