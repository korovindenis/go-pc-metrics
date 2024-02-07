package ip

import (
	"log"
	"net"
)

// GetOutbound returns the local IP address of the machine that connects to the specified server IP.
//
// net.IP
func GetOutbound() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

// CheckInSubnet checks if the given IP is in the specified subnet.
//
// Parameters:
// - ip string: the IP address to check
// - subnet string: the subnet to check against
// Return type(s):
// - bool: true if the IP is in the subnet, false otherwise
// - error: an error if the subnet parsing fails
func CheckInSubnet(ip string, subnet string) (bool, error) {
	_, subnetCIDR, err := net.ParseCIDR(subnet)
	if err != nil {
		return false, err
	}
	ipCIDR := net.ParseIP(ip)

	return subnetCIDR.Contains(ipCIDR), nil
}
