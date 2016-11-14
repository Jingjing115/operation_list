package parse_op

type MonitorOp struct {

}

func NewMonitorOp() *MonitorOp {
	return &MonitorOp{}
}

func (p *MonitorOp) Process(data interface{}) interface{} {
	if data == nil {
		return data
	}
	// TODO 处理监控层
	return data
}