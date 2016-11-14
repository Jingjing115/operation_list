package parse_op

import (
	"log"
)

type ParseOp struct {
	pipelines []Pipeline
}

func NewParseOp() *ParseOp {
	return &ParseOp{
		pipelines: []Pipeline{
			NewFilterOp(),
			NewMonitorOp(),
		},
	}
}

func (p *ParseOp) Connect() error {
	return nil
}

func (p *ParseOp) Start() error {
	log.Println("ParseOp Start...")
	return nil
}

func (p *ParseOp) Delivery(data interface{}) interface{} {
	if data == nil {
		return data
	}

	for _, line := range p.pipelines {
		data = line.Process(data)
	}
	return data
}
