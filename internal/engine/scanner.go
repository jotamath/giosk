package engine

import (
	"fmt"
	"net"
	"time"

	"giosk/internal/models"
)

func Worker(ips []string, ports <-chan int, results chan<- models.ScanResults, timeout time.Duration) {
	for port := range ports {
		for _, ip := range ips {
			results <- connect(ip, port, timeout)
		}
	}
}

func connect(ip string, port int, timeout time.Duration) models.ScanResults {
	addr := net.JoinHostPort(ip, fmt.Sprintf("%d", port))
	res := models.ScanResults{IP: ip, Port: port, Status: "closed", TimeStamp: time.Now().Unix()}

	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return res
	}
	defer conn.Close()

	res.Status = "open"

	conn.SetReadDeadline(time.Now().Add(time.Second * 2))
	buff := make([]byte, 256)
	n, _ := conn.Read(buff)
	if n > 0 {
		res.Banner = string(buff[:n])
	}

	return res
}
