package config

import (
	"testing"
	"fmt"
)

func TestNewConfig(t *testing.T) {
	conf := NewConfig("./config.yml")

	fmt.Println(conf)
}
