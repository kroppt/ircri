package conn

// Logger is used to log messages
type Logger interface {
	Log(msg string)
}

// ChannelLogger implements logging through channels
type ChannelLogger chan string

// Log sends the log message through the channel
func (cl ChannelLogger) Log(msg string) {
	cl <- msg
}

// NoopLogger does nothing with messages that get logged
type NoopLogger struct {
}

// Log does nothing
func (nl NoopLogger) Log(msg string) {
}
