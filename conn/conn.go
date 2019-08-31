package conn

import (
	"bufio"
	"fmt"
	"net"

	"github.com/pkg/errors"
)

// Client stores information about a client connection
type Client struct {
	tcp *net.TCPConn
}

// Connect creates a client connection with a server
func Connect(ip net.IP, port int) (*Client, error) {
	server := &net.TCPAddr{
		IP:   ip,
		Port: port,
	}
	conn, err := net.DialTCP("tcp", nil, server)
	if err != nil {
		return nil, errors.Wrap(err, "DialTCP failed")
	}
	fmt.Println("Connection established.")
	msg := "CAP LS 302\r\nNICK guest\r\nUSER guest 0 * :guest\r\n"
	fmt.Printf("Sending message '%s'...\n", msg)
	n, err := conn.Write([]byte(msg))
	if err != nil {
		return nil, errors.Wrap(err, "TCPConn.Write failed")
	}
	fmt.Printf("Wrote %v bytes.\n", n)
	fmt.Println("Reading message...")
	resp, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return nil, errors.Wrap(err, "TCPConn.Read failed")
	}
	fmt.Printf("Read message: '%s'.\n", resp)
	return &Client{conn}, nil
}

// Disconnect a client from its connected server
func (cl *Client) Disconnect() error {
	err := cl.tcp.Close()
	if err != nil {
		return errors.Wrap(err, "TCPConn.Close failed")
	}
	return nil
}
