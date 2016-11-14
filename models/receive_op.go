package models

import (
	"time"
	"encoding/json"
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

type ReceiveOp struct {
	Encrypted    bool   `json:"encrypted"`
	WEncrypted   bool   `json:"w_encrypted"`
	DeviceAddr   uint32 `json:"device_addr"`
	Op           uint16 `json:"op"`
	Params       string `json:"params"`
	UserKeyIndex int    `json:"user_key_index"`
	SrcCost      int    `json:"src_cost"`
	SrcSeq       int    `json:"src_seq"`
	Version      uint32 `json:"version"`
	TeleportAddr uint32 `json:"teleport_addr"`

	// 由生产者设置,单位 纳秒
	Timestamp    time.Time `json:"-"`
}

func (op *ReceiveOp) Equal(ano *ReceiveOp) bool {
	return op.Encrypted == ano.Encrypted &&
		op.WEncrypted == ano.WEncrypted &&
		op.DeviceAddr == ano.DeviceAddr &&
		op.Op == ano.Op &&
		op.Params == ano.Params &&
		op.UserKeyIndex == ano.UserKeyIndex &&
		op.SrcCost == ano.SrcCost &&
		op.SrcSeq == ano.SrcSeq &&
		op.Version == ano.Version &&
		op.TeleportAddr == ano.TeleportAddr
}

func (op *ReceiveOp) String() string {
	data, err := json.Marshal(op)
	if err != nil {
		return "{}"
	}
	return string(data)
}

func (op *ReceiveOp) Direction() bool {
	return true
}

func (op *ReceiveOp) Data() string {
	return op.String()
}

func (op *ReceiveOp) OpCode() string {
	return string(op.littleToBig() >> 8) + string(op.littleToBig() & 0xFF)
}

func (op *ReceiveOp) littleToBig() uint16 {
	code := make([]byte, 2)
	binary.LittleEndian.PutUint16(code, op.Op)
	return binary.BigEndian.Uint16(code)
}

func (op *ReceiveOp) ToCode() string {
	code := make([]byte, 2)
	binary.LittleEndian.PutUint16(code, op.Op)
	return hex.EncodeToString(code)
}

func (op *ReceiveOp) ToSnippet() string {
	return fmt.Sprintf("(to_timestamp(%s), 'null', %d, %d, '%s', '%s', TRUE, TRUE)",
		fmt.Sprintf("%d.%d", op.Timestamp.Unix(),
			uint64(op.Timestamp.UnixNano())-uint64(op.Timestamp.Unix())*10e8),
		op.DeviceAddr,
		op.TeleportAddr,
		op.ToCode(),
		op.Params,
	)
}


func GetTestReceiveOp() *ReceiveOp {
	return &ReceiveOp{
		Encrypted:    true,
		WEncrypted:   true,
		DeviceAddr:   12345,
		Op:           0x726f,
		Params:       "11223344",
		UserKeyIndex: 1,
		SrcCost:      2,
		SrcSeq:       3,
		Version:      4,
		TeleportAddr: 54321,

		Timestamp: time.Now(),
	}
}