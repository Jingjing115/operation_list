package conn

import (
	"github.com/garyburd/redigo/redis"
	"strconv"
	"log"
	"fmt"
	"time"
)

type RedisConn struct {
	host      string
	port      string
	password  string
	database  string

	conn      redis.Conn
	pubConn   *redis.PubSubConn
	connErr   chan redis.Error

	reconnect chan bool
}

func NewRedisConn(host, port, password, database string) *RedisConn {
	return &RedisConn{
		host: host, port: port,
		password: password,
		database: database,
		connErr:make(chan redis.Error),
		reconnect:make(chan bool, 1),
	}
}

func (c *RedisConn) Connect() error {
	c.conn = connectToRedis(c.host, c.port, c.password, c.database)
	c.pubConn = &redis.PubSubConn{Conn: connectToRedis(c.host, c.port, c.password, c.database)}
	c.recovery()
	return nil
}

func connectToRedis(host, port, password, database string) redis.Conn {
	db, err := strconv.Atoi(database)
	if err != nil {
		log.Fatalln(err)
	}
	opPass := redis.DialPassword(password)
	opDb := redis.DialDatabase(db)
	for {
		conn, err := redis.Dial("tcp", fmt.Sprintf("%s:%s", host, port), opPass, opDb)
		if err == nil {
			return conn
		}
		log.Println(err)
		log.Printf("Trying to reconnect to Redis at %s:%s\n", host, port)
		time.Sleep(1 * time.Second)
	}
}

func (c *RedisConn) recovery() {
	go func() {
		for {
			<-c.connErr
			log.Printf("Connecting to redis %s:%s\n", c.host, c.port)
			c.conn = connectToRedis(c.host, c.port, c.password, c.database)
			c.pubConn = &redis.PubSubConn{Conn: connectToRedis(c.host, c.port, c.password, c.database)}
			c.reconnect <- true
		}
	}()
}


func (c *RedisConn) PubSubConn() *redis.PubSubConn {
	return c.pubConn
}

func (c *RedisConn) Receive() interface{} {
	return c.pubConn.Receive()
}

func (c *RedisConn) Subscribe(channels ...string) error {
	var err error
	for _, ch := range channels {
		if err = c.pubConn.Subscribe(ch); err != nil {
			return err
		}
	}
	return nil
}

func (c *RedisConn) Publish(channel string, data interface{}) error {
	_, err := c.conn.Do("PUBLISH", channel, data)
	return err
}

func (c *RedisConn) ConnErr(err string) {
	c.connErr <- redis.Error(err)
}

func (c *RedisConn) Reconnect() chan bool {
	return c.reconnect
}

func (c *RedisConn) Adapter() string {
	return "redis"
}
