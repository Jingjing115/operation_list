package parse_op

type Pipeline interface {
	Process(data interface{}) interface{}
}
