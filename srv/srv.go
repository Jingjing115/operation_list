package srv

import (
	"log"
	"parse_op/config"
	"parse_op/data_receive"
	"parse_op/msg_distribute"
	"parse_op/parse_op"
)

type Srv struct {
	modules []Modules
}

func NewSrv(configFile string) *Srv {
	conf := config.NewConfig(configFile)
	dc := conf.DataReceiveConf
	mc := conf.MsgDistributeConf
	return &Srv{
		modules: []Modules{
			data_receive.NewDataReceive(dc.Producer, dc.Host, dc.Port, dc.Username, dc.Password, dc.Database, dc.Rch, dc.Pch),
			parse_op.NewParseOp(),
			msg_distribute.NewMsgDistribute(mc.Consumer, mc.Host, mc.Port, mc.Username, mc.Password, mc.Database, mc.WhiteList),
		},
	}
}

func (s *Srv) Start() error {
	var err error
	for _, module := range s.modules {
		// 1. 先连接
		if err = module.Connect(); err != nil {
			log.Fatalf("Server connecting error: %s\n", err)
			return err
		}
		// 2. 后开始
		if err = module.Start(); err != nil {
			log.Fatalf("Server start error: %s\n", err)
			return err
		}
	}
	return nil
}

// 死循环处理数据流
func (s *Srv) Process() {
	var data interface{}
	for {
		for _, module := range s.modules {
			data = module.Delivery(data)
			if data == nil {
				break
			}
		}
	}
}
