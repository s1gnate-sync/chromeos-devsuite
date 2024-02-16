package ssh

import (
	"fmt"
	"github.com/creack/pty"
	"github.com/gliderlabs/ssh"
	"io"
	"net"
	"os/exec"
)

type CmdProvider func(name string, args []string) *exec.Cmd
type CmdProviders map[string]CmdProvider

func Run(address string, hostKeyFile string, providers CmdProviders) error {
	sessionRequestCallback := func(s ssh.Session, requestType string) bool {
		if providers[s.User()] == nil {
			s.Exit(1)
			return false
		}

		return true
	}

	server := &ssh.Server{
		PtyCallback:                   ptyCallback,
		PublicKeyHandler:              publicKeyHandler,
		PasswordHandler:               passwordHandler,
		LocalPortForwardingCallback:   localPortForwardingCallback,
		ReversePortForwardingCallback: reversePortForwardingCallback,
		SessionRequestCallback:        sessionRequestCallback,
		KeyboardInteractiveHandler:    keyboardInteractiveHandler,
	}

	conn, err := net.Listen("tcp", address)
	if err != nil {
		return err

	}
	signer, err := provideHostKey(hostKeyFile)
	if err != nil {
		return err
	}

	server.AddHostKey(signer)
	server.Handler = func(s ssh.Session) {
		ptyReq, winCh, doPty := s.Pty()
		cmd := providers[s.User()](s.User(), s.Command())

		if cmd == nil {
			return
		}

		if doPty {
			cmd.Env = append(cmd.Env, fmt.Sprintf("TERM=%s", ptyReq.Term))

			f, err := pty.Start(cmd)
			if err != nil {
				return
			}

			go func() {
				for win := range winCh {
					setWinsize(f, win.Width, win.Height)
				}
			}()

			go func() {
				io.Copy(f, s)
			}()

			io.Copy(s, f)
			cmd.Wait()
		} else {
			cmd.Stdout = s
			cmd.Stderr = s.Stderr()

			if err := cmd.Run(); err != nil {
				s.Exit(1)
			}
		}
	}

	if err := server.Serve(conn); err != nil {
		return err
	}

	return nil
}
