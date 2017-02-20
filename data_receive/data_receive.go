package data_receive

import (
	"github.com/garyburd/redigo/redis"
	"github.com/streadway/amqp"
	"log"
	"time"
	"parse_op/models"
	"encoding/json"
	"parse_op/conn"
)

type DataReceive struct {
	producer  string // rabbit, redis
	host      string
	port        string
	username  string
	password  string
	database  string

	rch       string
	pch       string

	conn      conn.Conn
	connErr   chan bool
	reconnect chan bool

	data      chan interface{}
}

func NewDataReceive(producer, host, port, username, password, database, rch, pch string) *DataReceive {
	d := &DataReceive{
		producer: producer, host: host, port: port,
		username: username, password: password, database: database,
		rch:rch, pch:pch,
		reconnect: make(chan bool, 1),
		data: make(chan interface{}, 1),
	}
	d.initConn()
	return d
}

func (d *DataReceive) initConn() {
	switch d.producer {
	case "redis":
		d.conn = conn.NewRedisConn(d.host, d.port, d.password, d.database)
	case "rabbit":
		d.conn = conn.NewRabbitConn(d.host, d.port, d.username, d.password)
	}
}

func (d *DataReceive) Connect() (err error) {
	if err = d.conn.Connect(); err != nil {
		return
	}
	d.recovery()
	return
}

func (d *DataReceive) startRedis() (err error) {
	if err = d.conn.(*conn.RedisConn).Subscribe(d.rch, d.pch); err != nil {
		return err
	}
	go d.receiveRedisData()

	return nil
}

func (d *DataReceive) receiveRedisData() {
	c := d.conn.(*conn.RedisConn)
	for {
		data := c.Receive()

		// 读出错误，redis断线，重连
		if err, ok := data.(error); ok {
			c.ConnErr("redis read fatal.")
			log.Printf("redis fatal: %s\n", err.Error())
			return
		}
		if n, ok := data.(redis.Message); ok {
			handleTime := time.Now()
			if n.Channel == d.rch {
				op := new(models.ReceiveOp)
				err := json.Unmarshal(n.Data, op)
				op.Timestamp = handleTime
				if err == nil {
					d.data <- op
				}
			} else if n.Channel == d.pch {
				ops := make([]models.PushOp, 0, 1)
				err := json.Unmarshal(n.Data, &ops)
				if err == nil {
					for _, op := range ops {
						op.Timestamp = handleTime
						d.data <- &op
					}
				}
			}
		}
	}
}

func (d *DataReceive) receiveRabbitData(rmsg, pmsg <-chan amqp.Delivery) {
	for {
		handleTime := time.Now()
		select {
		case msg := <-rmsg:
			op := new(models.ReceiveOp)
			err := json.Unmarshal(msg.Body, op)
			op.Timestamp = handleTime
			if err == nil {
				d.data <- op
			}
		case msg := <-pmsg:
			ops := make([]models.PushOp, 0, 1)
			err := json.Unmarshal(msg.Body, &ops)
			if err == nil {
				for _, op := range ops {
					op.Timestamp = handleTime
					d.data <- &op
				}
			}
		case re := <-d.reconnect:
			log.Println("recover rabbit: ", re)
			return
		}
	}
}

func (d *DataReceive) startRabbit() (err error) {
	c := d.conn.(*conn.RabbitConn)

	rmsg, err := c.Subscribe(d.rch, "")
	if err != nil {
		return
	}

	pmsg, err := c.Subscribe(d.pch, "")
	if err != nil {
		return
	}
	go d.receiveRabbitData(rmsg, pmsg)

	return nil
}

func (d *DataReceive) recovery() {
	go func() {
		for {
			<-d.conn.Reconnect()
			d.reconnect <- true
			if err := d.Start(); err != nil {
				log.Println(err)
				return
			}
		}
	}()
}

func (d *DataReceive) Start() (err error) {
	log.Println("DataReceive Start...")
	d.Connect()
	if d.producer == "rabbit" {
		err = d.startRabbit()
	} else if d.producer == "redis" {
		err = d.startRedis()
	}
	return
}

func (d *DataReceive) Delivery(data interface{}) interface{} {
	return <-d.data
}
