package conn

import "testing"

func TestChannelLoggerBasic(t *testing.T) {
	log := make(chan string, 1)
	cl := ChannelLogger(log)
	msg := "Test log message"
	cl.Log(msg)
	recvmsg := <-log
	if msg != recvmsg {
		t.Errorf("Expected logged message '%s', but got '%s'", msg, recvmsg)
	}
}
