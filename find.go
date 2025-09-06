package server_with_register

import (
	"fmt"
	"net"
)

const (
	portStart = 12300
	portEnd   = 12500
)

func findAvailablePort() (int64, error) {
	// build a consul client
	for port := portStart; port < portEnd; port++ {
		if IsPortAvailable(port) {
			return int64(port), nil
		}
	}
	return 0, fmt.Errorf("no available port")
}

func IsPortAvailable(port int) bool {
	address := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return false // 端口被占用
	}
	defer listener.Close()
	return true // 端口可用
}
