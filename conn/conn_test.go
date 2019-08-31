package conn

import (
	"net"
	"testing"
)

func TestConnectBasic(t *testing.T) {
	cl, err := Connect(net.IPv4(91, 236, 182, 1), 6667)
	if err != nil {
		t.Error(err)
	}
	if cl == nil {
		t.Error("returned conn.Client is nil")
	}
}

func TestDisconnect(t *testing.T) {
	cl, _ := Connect(net.IPv4(91, 236, 182, 1), 6667)
	t.Run("TestDisconnect", func(t *testing.T) {
		cl.Disconnect()
	})
}
