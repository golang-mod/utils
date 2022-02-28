package ip

import (
	"net"
	"net/http"
	"strings"
)

func GetRealIp(r *http.Request) string {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		ip = r.RemoteAddr
	}
	if ip != "127.0.0.1" {
		return ip
	}
	// Check if behide nginx or apache
	xRealIP := r.Header.Get("X-Real-Ip")
	xForwardedFor := r.Header.Get("X-Forwarded-For")

	for _, address := range strings.Split(xForwardedFor, ",") {
		address = strings.TrimSpace(address)
		if address != "" {
			return address
		}
	}

	if xRealIP != "" {
		return xRealIP
	}
	return ip
}
