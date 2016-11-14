package conn

import (
	"github.com/streadway/amqp"
	"log"
	"time"
	"fmt"

	"strings"
)

type RabbitConn struct {
	host      string
	port      string
	username  string
	password  string

	conn      *amqp.Connection
	connErr   chan *amqp.Error

	reconnect chan bool
}

func NewRabbitConn(host, port, username, password string) *RabbitConn {
	return &RabbitConn{
		host: host, port: port,
		username: username, password: password,
		connErr: make(chan *amqp.Error),
		reconnect: make(chan bool, 1),
	}
}

func (c *RabbitConn) Connect() error {
	c.conn = connectToRabbit(c.url())
	c.conn.NotifyClose(c.connErr)
	c.recovery()
	return nil
}

// 每秒尝试连接一次rabbit
func connectToRabbit(url string) *amqp.Connection {
	for {
		conn, err := amqp.Dial(url)
		if err == nil {
			return conn
		}
		log.Println(err)
		log.Printf("Trying to reconnect to RabbitMQ at %s\n", url)
		time.Sleep(1 * time.Second)
	}
}

func (c *RabbitConn) url() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s/", c.username, c.password, c.host, c.port)
}

func (c *RabbitConn) recovery() {
	go func() {
		for {
			<-c.connErr
			log.Printf("Connecting to %s\n", c.url())
			c.conn = connectToRabbit(c.url())
			c.connErr = make(chan *amqp.Error)
			c.conn.NotifyClose(c.connErr)
			c.reconnect <- true
		}
	}()
}

func (c *RabbitConn) Conn() *amqp.Connection {
	return c.conn
}


func (c *RabbitConn) Publish(queue string, data interface{}) error {
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}
	var body []byte
	switch data.(type) {
	case string:
		body = []byte(data.(string))
	case []byte:
		body = data.([]byte)
	case []string:
		stringByte := "[" + strings.Join(data.([]string), ",") + "]"
		body = []byte(stringByte)
	}

	err = ch.Publish(
		"",
		queue,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType: "text/plain",
			Body: body,
		},
	)
	return err
}

func (c *RabbitConn) Subscribe(queue, consumer string) (<-chan amqp.Delivery, error) {
	ch, err := c.conn.Channel()
	if err != nil {
		return nil, err
	}
	q, err := ch.QueueDeclare(
		queue, // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil, // arguments
	)
	if err != nil {
		return nil, err
	}
	return ch.Consume(
		q.Name,
		consumer,
		true, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil, // args
	)
}

func (c *RabbitConn) BindQueue(queue string) (*amqp.Channel, error) {
	ch, err := c.conn.Channel()
	if err != nil {
		return nil, err
	}
	_, err = ch.QueueDeclare(
		queue, // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil, // arguments
	)
	if err != nil {
		return nil, err
	}
	return ch, nil
}

func (c *RabbitConn) Reconnect() chan bool {
	return c.reconnect
}

func (c *RabbitConn) Adapter() string {
	return "rabbit"
}