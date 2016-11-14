package models

type OP interface {
	OpCode() string
	Direction() bool
	Data() string
}
