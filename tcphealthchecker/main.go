package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

func main() {
	checker := NewTCPChecker(net.ParseIP("127.0.0.1"), 4321, 10)
	checker.Timeout = 1 * time.Second

	logOutput := log.Writer()
	result := checker.CheckWithRetries(5, 4*time.Second, logOutput)

	println("Result: ", result.Message)
}

type Target struct {
	IP      net.IP
	Host    net.IP
	Port    int
	Packets int
}

type TCPChecker struct {
	Target
	Timeout time.Duration
}

type Result struct {
	Success bool
	Message string
}

func NewTCPChecker(ip net.IP, port int, packets int) *TCPChecker {
	return &TCPChecker{
		Target: Target{
			IP:      ip,
			Port:    port,
			Packets: packets,
		},
	}
}

func (hc *TCPChecker) addr() string {
	return fmt.Sprintf("%s:%d", hc.IP.String(), hc.Port)
}

func (hc *TCPChecker) Check(timeout time.Duration) *Result {
	conn, err := net.DialTimeout("tcp", hc.addr(), timeout)

	if err != nil {
		return &Result{Success: false, Message: fmt.Sprintf("Failed to connect: %v", err)}
	}
	defer conn.Close() // graceful close connection

	// here you can do a bunch of checks on the network, like tls, packets lost rate and many more things

	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Printf("conn.Read() error: %v", err)
			}
			break
		}

		// print message received
		log.Printf("Received: %q", buf[:n])
	}

	return &Result{Success: true, Message: "Connected succesfully"}
}

func (hc *TCPChecker) CheckWithRetries(retries int, retryDelay time.Duration, logOutput io.Writer) *Result {
	var result *Result

	for i := 0; i < retries; i++ {
		start := time.Now()
		result = hc.Check(hc.Timeout)
		duration := time.Since(start)

		logOutput.Write([]byte(fmt.Sprintf("Health Check Attempt %d - Success: %v, Latency: %v, MSG: %s\n", i+1, result.Success, duration, result.Message)))

		if result.Success {
			return result
		}

		// if not succesful, try again
		time.Sleep(retryDelay)
	}

	return result
}
