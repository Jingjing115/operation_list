package msg_distribute

import (
	"testing"
	"parse_op/models"
	"parse_op/conn"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"encoding/json"
)

func TestNewMsgDistribute(t *testing.T) {
	consumer := "redis"
	host := "127.0.0.1"
	port := "6379"
	username := ""
	password := ""
	database := "0"
	list := "ro"
	var err error

	d := NewMsgDistribute(consumer, host, port, username, password, database, list)
	d.Connect()
	rop := models.GetTestReceiveOp()

	c := d.conn.(*conn.RedisConn)
	if err = c.Subscribe(list); err != nil {
		t.Fatal(err)
	}
	if n, ok := c.Receive().(redis.Subscription); ok {
		fmt.Println(n)
	} else {
		t.Fatal("error subscribe.")
	}

	d.publishToRedis(rop)

	if n, ok := c.Receive().(redis.Message); ok {
		if n.Channel == list {
			op := new(models.ReceiveOp)
			if err = json.Unmarshal(n.Data, &op); err != nil {
				t.Fatal(err)
			}
			if !op.Equal(rop) {
				t.Fatalf("error op expect:%s, got:%s.\n", rop, op)
			}
		}
	}

}
