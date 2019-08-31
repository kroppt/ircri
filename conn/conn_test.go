package conn

import "testing"

func TestConnect(t *testing.T) {
	t.Run("TestConnect", func(t *testing.T) {
		Connect()
	})
}

func TestDisconnect(t *testing.T) {
	cl := Connect()
	t.Run("TestDisconnect", func(t *testing.T) {
		cl.Disconnect()
	})
}
