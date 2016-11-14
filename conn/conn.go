package conn

type Conn interface {
	Adapter() string
	Connect() error
	Reconnect() chan bool
}
