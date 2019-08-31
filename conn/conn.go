package conn

import (
	"bufio"
	"fmt"
	"net"

	"github.com/pkg/errors"
)

// ClientOptions configures the client on creation
type ClientOptions struct {
	Logger Logger
}

// NewClientOptions returns a ClientOptions with valid values filled in
func NewClientOptions(opts ClientOptions) *ClientOptions {
	if opts.Logger == nil {
		opts.Logger = NoopLogger{}
	}
	return &opts
}

// Client stores information about a client connection
type Client struct {
	tcp  *net.TCPConn
	opts *ClientOptions
}

// Connect creates a client connection with a server
func Connect(ip net.IP, port int, opts *ClientOptions) (*Client, error) {
	server := &net.TCPAddr{
		IP:   ip,
		Port: port,
	}
	conn, err := net.DialTCP("tcp", nil, server)
	if err != nil {
		return nil, errors.Wrap(err, "DialTCP failed")
	}
	opts.Logger.Log("Connection established.")
	msg := "CAP LS 302\r\nNICK guest\r\nUSER guest 0 * :guest\r\n"
	opts.Logger.Log(fmt.Sprintf("Sending message '%s'...\n", msg))
	n, err := conn.Write([]byte(msg))
	if err != nil {
		return nil, errors.Wrap(err, "TCPConn.Write failed")
	}
	opts.Logger.Log(fmt.Sprintf("Wrote %v bytes.\n", n))
	opts.Logger.Log("Reading message...")
	resp, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return nil, errors.Wrap(err, "TCPConn.Read failed")
	}
	opts.Logger.Log(fmt.Sprintf("Read message: '%s'.\n", resp))
	return &Client{conn, opts}, nil
}

// Disconnect a client from its connected server
func (cl *Client) Disconnect() error {
	err := cl.tcp.Close()
	if err != nil {
		return errors.Wrap(err, "TCPConn.Close failed")
	}
	return nil
}
