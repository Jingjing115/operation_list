package models

import "testing"

func TestReceiveOp_OpCode(t *testing.T) {
	op := GetTestReceiveOp()

	code := "rb"
	op.Op = 0x7262
	if op.OpCode() != code {
		t.Fatalf("opcode parse error, expect: %s, got: %s.", code, op.OpCode())
	}

	code = "qi"
	op.Op = 0x7169
	if op.OpCode() != code {
		t.Fatalf("opcode parse error, expect: %s, got: %s.", code, op.OpCode())
	}

	code = "in"
	op.Op = 0x696e
	if op.OpCode() != code {
		t.Fatalf("opcode parse error, expect: %s, got: %s.", code, op.OpCode())
	}

	code = "qt"
	op.Op = 0x7174
	if op.OpCode() != code {
		t.Fatalf("opcode parse error, expect: %s, got: %s.", code, op.OpCode())
	}

	code = "at"
	op.Op = 0x6174
	if op.OpCode() != code {
		t.Fatalf("opcode parse error, expect: %s, got: %s.", code, op.OpCode())
	}

	code = "po"
	op.Op = 0x706f
	if op.OpCode() != code {
		t.Fatalf("opcode parse error, expect: %s, got: %s.", code, op.OpCode())
	}
}
