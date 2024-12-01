package isunippets

import (
	"net"
	"os"
)

type CreateUnixDomainSocketParams struct {
	Address string
}

func CreateUnixDomainSocket(address string) (net.Listener, error) {
	if stat, err := os.Stat(address); err == nil {
		if stat.Mode()&os.ModeSocket == 0 {
			return nil, os.ErrExist
		}
		err := os.Remove(address)
		if err != nil {
			return nil, err
		}
	}
	l, err := net.Listen("unix", address)
	if err != nil {
		return nil, err
	}
	err = os.Chmod(address, 0777)
	if err != nil {
		return nil, err
	}
	return l, nil
}

func GetUnusedPort() (int, error) {
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}
