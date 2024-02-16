package ssh

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/pem"
	gossh "golang.org/x/crypto/ssh"
	"io/ioutil"
	"os"
)

func provideHostKey(filename string) (gossh.Signer, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err == nil {
		signer, err := gossh.ParsePrivateKey(bytes)
		if err != nil {
			return nil, err
		}

		return signer, nil
	}

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	pemBytes, err := gossh.MarshalPrivateKey(key, "")
	if err != nil {
		return nil, err
	}

	bytes = pem.EncodeToMemory(pemBytes)
	if err := os.WriteFile(filename, bytes, 0700); err != nil {
		return nil, err
	}

	signer, err := gossh.ParsePrivateKey(bytes)
	if err != nil {
		return nil, err
	}

	return signer, nil
}
