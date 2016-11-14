package models

import (
	"testing"
	"fmt"
)

func TestPushOp_ToCode(t *testing.T) {
	op := GetTestPushOp()
	fmt.Println(op.OpCode())
}