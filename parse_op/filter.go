package parse_op

import (
	"parse_op/models"
	"math"
	"errors"
)

type FilterOp struct {
	filters [](func(op interface{}) error)
}

// 过滤无效的设备IP地址
func filterInvalidDeviceIp(op interface{}) error {
	switch x := op.(type) {
	case *models.ReceiveOp:
		if x.DeviceAddr > math.MaxInt32 {
			return errors.New("Device ip too large!")
		}
	case *models.PushOp:
		if x.Command.DeviceAddr > math.MaxInt32 {
			return errors.New("Device ip too large!")
		}
	}
	return nil
}

func NewFilterOp() *FilterOp {
	return &FilterOp{
		filters: [](func(op interface{}) error) {
			filterInvalidDeviceIp,
		},
	}
}

// 实现pipeline的接口
func (p *FilterOp) Process(data interface{}) interface{} {
	if data == nil {
		return data
	}
	// 使用所有的过滤器
	for _, filter := range p.filters {
		if err := filter(data); err != nil {
			return nil
		}
	}
	return data
}

