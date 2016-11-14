package parse_op

import (
	"testing"
	"parse_op/models"
)

// 过滤无效
func TestFilterOp_Process(t *testing.T) {
	p := NewFilterOp()
	op := models.GetTestReceiveOp()
	op.DeviceAddr = 0xFFFFFFFF

	newOp := p.Process(op)

	if newOp != nil {
		t.Fatalf("Filter missing, got: %s", newOp)
	}
}

// 保留有效
func TestFilterOp_Process2(t *testing.T) {
	p := NewFilterOp()
	op := models.GetTestReceiveOp()
	op.DeviceAddr = 0xFFFF

	newOp := p.Process(op)

	if newOp == nil {
		t.Fatalf("Filter missing, got: %s", newOp)
	}

	if _op, ok := newOp.(*models.ReceiveOp); ok {
		if !_op.Equal(op) {
			t.Fatalf("Filter fatal, expect %s, got %s.", op.String(), _op.String())
		}
	} else {
		t.Fatal("Filter changed op's type")
	}
}
