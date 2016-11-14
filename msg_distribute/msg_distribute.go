package msg_distribute

import (
	"github.com/streadway/amqp"
	"log"
	"parse_op/conn"
	"parse_op/models"
	"strings"
	"fmt"
)

type MsgDistribute struct {
	consumer  string // rabbit, redis
	host      string
	port      string
	username  string
	password  string
	database  string

	whiteList []string

	conn      conn.Conn
	connErr   chan bool
	reconnect chan bool

	data      chan interface{}
}

func NewMsgDistribute(consumer, host, port, username, password, database, list string) *MsgDistribute {
	msg := &MsgDistribute{
		consumer: consumer, host: host, port: port,
		username: username, password: password, database: database,
		reconnect: make(chan bool, 1),
		data: make(chan interface{}, 1),
	}
	msg.whiteList = strings.Split(list, ",")
	msg.initConn()
	return msg
}

func (d *MsgDistribute) initConn() {
	switch d.consumer {
	case "redis":
		d.conn = conn.NewRedisConn(d.host, d.port, d.password, d.database)
	case "rabbit":
		d.conn = conn.NewRabbitConn(d.host, d.port, d.username, d.password)
	}
}

func (d *MsgDistribute) Connect() (err error) {
	if err = d.conn.Connect(); err != nil {
		return
	}
	d.recovery()
	return
}


func (d *MsgDistribute) publishToRabbit(data interface{}) {
	op := data.(models.OP)
	c := d.conn.(*conn.RabbitConn)
	ch, err := c.Conn().Channel()
	defer ch.Close()
	if err != nil {
		log.Println(err)
		return
	}
	if err = ch.ExchangeDeclare("parse_op", "direct", false, false, false, false, nil); err != nil {
		log.Println(err)
		return
	}

	err = ch.Publish("parse_op", op.OpCode(), false, false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body: []byte(op.Data()),
		},
	)
	if err != nil {
		log.Println(err)
		return
	}
}

func (d *MsgDistribute) publishToRedis(data interface{}) {
	op := data.(models.OP)
	for _, list := range d.whiteList {
		if op.OpCode() == list {
			fmt.Printf("send msg %s to %s\n", op.Data(), "op_"+ op.OpCode())
			if err := d.conn.(*conn.RedisConn).Publish("op_" + op.OpCode(), op.Data()); err != nil {
				log.Println(err)
			}
		}
	}
}


func (d *MsgDistribute) Delivery(data interface{}) interface{} {
	if data == nil {
		return data
	}

	d.data <- data

	return nil
}

func (d *MsgDistribute) recovery() {
	go func() {
		for {
			<-d.conn.Reconnect()
			d.reconnect <- true
			if err := d.Connect(); err != nil {
				log.Println(err)
			}
		}
	}()
}

func (d *MsgDistribute) Start() (err error) {
	log.Println("MsgDistribute Start...")
	go func() {
		for {
			data := <-d.data
			if d.consumer == "rabbit" {
				d.publishToRabbit(data)
			} else if d.consumer == "redis" {
				d.publishToRedis(data)
			}
		}
	}()
	return nil
}