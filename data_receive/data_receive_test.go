package data_receive

import (
	"testing"
	"parse_op/conn"
	"parse_op/models"

)

func TestNewDataReceiveRedis(t *testing.T) {
	producer := "redis"
	host := "127.0.0.1"
	port := "6379"
	username := ""
	password := ""
	database := "0"
	rch := "tcp_receive"
	pch := "tcp_push"
	var err error

	d := NewDataReceive(producer, host, port, username, password, database, rch, pch)
	if err = d.Connect(); err != nil {
		t.Fatal(err)
	}

	if err = d.Start(); err != nil {
		t.Fatal(err)
	}

	rop := models.GetTestReceiveOp()
	if err := d.conn.(*conn.RedisConn).Publish(rch, rop.String()); err != nil {
		t.Fatal(err)
	}

	data := d.Delivery(nil)
	if op, ok := data.(*models.ReceiveOp); ok {
		if !op.Equal(rop) {
			t.Fatalf("redis err, expect %s, got %s", rop, op)
		}
	} else {
		t.Fatal("wrong op type.")
	}

	pop := models.GetTestPushOp()
	if err = d.conn.(*conn.RedisConn).Publish(pch, []string{pop.String()}); err != nil {
		t.Fatal(err)
	}
	data = d.Delivery(nil)
	if op, ok := data.(*models.PushOp); ok {
		if !op.Equal(pop) {
			t.Fatalf("redis err, expect %s, got %s", pop, op)
		}
	} else {
		t.Fatal("wrong op type.")
	}
}

func TestNewDataReceiveRabbit(t *testing.T) {
	producer := "rabbit"
	host := "127.0.0.1"
	port := "5672"
	username := "guest"
	password := "guest"
	database := "0"
	rch := "tcp_receive"
	pch := "tcp_push"
	var err error


	d := NewDataReceive(producer, host, port, username, password, database, rch, pch)

	if err = d.Connect(); err != nil {
		t.Fatal(err)
	}

	if err = d.Start(); err != nil {
		t.Fatal(err)
	}
	rop := models.GetTestReceiveOp()
	if err = d.conn.(*conn.RabbitConn).Publish(rch, rop.String()); err != nil {
		t.Fatal(err)
	}

	data := d.Delivery(nil)
	if op, ok := data.(*models.ReceiveOp); ok {
		if !op.Equal(rop) {
			t.Fatalf("redis err, expect %s, got %s", rop, op)
		}
	} else {
		t.Fatal("wrong op type.")
	}

	pop := models.GetTestPushOp()
	if err = d.conn.(*conn.RabbitConn).Publish(pch, []string{pop.String()}); err != nil {
		t.Fatal(err)
	}

	data = d.Delivery(nil)
	if op, ok := data.(*models.PushOp); ok {
		if !op.Equal(pop) {
			t.Fatalf("redis err, expect %s, got %s", pop, op)
		}
	} else {
		t.Fatal("wrong op type.")
	}
}

//func TestNewDataReceiveRabbit2(t *testing.T) {
//	producer := "rabbit"
//	host := "127.0.0.1"
//	port := "5672"
//	username := "guest"
//	password := "guest"
//	database := "0"
//	rch := "tcp_receive"
//	pch := "tcp_push"
//	var err error
//
//
//	d := NewDataReceive(producer, host, port, username, password, database, rch, pch)
//
//	if err = d.Connect(); err != nil {
//		t.Fatal(err)
//	}
//
//	if err = d.Start(); err != nil {
//		t.Fatal(err)
//	}
//	rop := models.GetTestReceiveOp()
//	if err = d.conn.(*conn.RabbitConn).Publish(rch, rop.String()); err != nil {
//		t.Fatal(err)
//	}
//
//	data := d.Delivery(nil)
//	if op, ok := data.(*models.ReceiveOp); ok {
//		if !op.Equal(rop) {
//			t.Fatalf("redis err, expect %s, got %s", rop, op)
//		}
//	} else {
//		t.Fatal("wrong op type.")
//	}
//	pop := models.GetTestPushOp()
//	go func() {
//		for {
//			if err = d.conn.(*conn.RabbitConn).Publish(pch, []string{pop.String()}); err != nil {
//				t.Fatal(err)
//			}
//			time.Sleep(1 * time.Second)
//		}
//	}()
//
//	go func() {
//		for {
//			data = d.Delivery(nil)
//			log.Println(data)
//			if op, ok := data.(*models.PushOp); ok {
//				if !op.Equal(pop) {
//					t.Fatalf("redis err, expect %s, got %s", pop, op)
//				}
//			} else {
//				t.Fatal("wrong op type.")
//			}
//			time.Sleep(1 * time.Second)
//		}
//	}()
//
//	//forever := make(chan bool)
//	//forever <- true
//
//}
