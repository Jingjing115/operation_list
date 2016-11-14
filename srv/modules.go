package srv

type Modules interface {
	Connect() error
	Delivery(data interface{}) interface{}
	Start() error
}
