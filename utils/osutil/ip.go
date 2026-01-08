package osutil

import (
	"errors"
	"net"
	"strings"
)

// GetNetIP retrieves the local IP address associated with the default network gateway.
// It establishes a UDP connection to google.com:80 to determine which local network interface is used for outbound traffic. If the machine has multiple interfaces, it returns the IP of the primary interface used for the route.
func GetDefaultNetIP() (ip string, err error) {
	conn, err := net.Dial("udp", "google.com:80")
	if err != nil {
		return "", errors.New("get net ip failed")
	}
	defer conn.Close()
	return strings.Split(conn.LocalAddr().String(), ":")[0], nil
}
