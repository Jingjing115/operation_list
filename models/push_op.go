package models

import (
	"time"
	"encoding/json"
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

type PushOpCommand struct {
	Encrypted  bool   `json:"encrypted"`
	WEncrypted bool   `json:"w_encrypted"`
	DeviceAddr uint32 `json:"device_addr"`
	Op         uint16 `json:"op"`
	Params     string `json:"params"`
}

type PushOp struct {
	Type         string `json:"type"`
	TeleportAddr uint32 `json:"teleport_addr"`
	Command      *PushOpCommand `json:"command"`

	Timestamp    time.Time `json:"-"`
}

func (opCmd *PushOpCommand) Equal(ano *PushOpCommand) bool {
	return opCmd.Encrypted == ano.Encrypted &&
		opCmd.WEncrypted == ano.WEncrypted &&
		opCmd.DeviceAddr == ano.DeviceAddr &&
		opCmd.Op == ano.Op &&
		opCmd.Params == ano.Params
}

func (op *PushOp) Equal(ano *PushOp) bool {
	return op.Type == ano.Type &&
		op.TeleportAddr == ano.TeleportAddr &&
		op.Command.Equal(ano.Command)
}

func (op *PushOp) String() string {
	data, err := json.Marshal(op)
	if err != nil {
		return "{}"
	}
	return string(data)
}

func (op *PushOp) Direction() bool {
	return false
}

func (op *PushOp) Data() string {
	return op.String()
}

func (op *PushOp) OpCode() string {
	return string(op.littleToBig() >> 8) + string(op.littleToBig() & 0xFF)
}

func (op *PushOp) littleToBig() uint16 {
	code := make([]byte, 2)
	binary.LittleEndian.PutUint16(code, op.Command.Op)
	return binary.BigEndian.Uint16(code)
}

func (op *PushOp) ToCode() string {
	code := make([]byte, 2)
	binary.LittleEndian.PutUint16(code, op.Command.Op)
	return hex.EncodeToString(code)
}

func (op *PushOp) ToSnippet() string {
	return fmt.Sprintf("(to_timestamp(%s), 'null', %d, %d, '%s', '%s', TRUE, TRUE)",
		fmt.Sprintf("%d.%d", op.Timestamp.Unix(),
			uint64(op.Timestamp.UnixNano()) - uint64(op.Timestamp.Unix()) * 10e8),
		op.Command.DeviceAddr,
		op.TeleportAddr,
		op.ToCode(),
		op.Command.Params,
	)
}

func GetTestPushOp() *PushOp {
	return &PushOp{
		Type:         "tcp",
		TeleportAddr: 54321,
		Command: &PushOpCommand{
			Encrypted:  true,
			WEncrypted: true,
			DeviceAddr: 12345,
			Op:         0x726f,
			Params:     "44332211",
		},
		Timestamp: time.Now(),
	}
}