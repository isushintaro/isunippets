package isunippets

import (
	"net"
)

type CreateUnixDomainSocketParams struct {
	Address string
}

func CreateUnixDomainSocket(address string) (net.Listener, error) {
	l, err := net.Listen("unix", address)
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
