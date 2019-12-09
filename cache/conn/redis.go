package conn

import (
	"github.com/garyburd/redigo/redis"
	"log"
	"time"
)

var (
	pool      *redis.Pool
	redisHost = "127.0.0.1:6379"
	redisPwd  = "123456"
)

func newRedisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     50,
		MaxActive:   30,
		IdleTimeout: 300 * time.Second,
		Dial: func() (conn redis.Conn, err error) {
			if conn, err = redis.Dial("tcp", redisHost); err != nil {
				log.Println(err.Error())
				return
			}
			if _, err = conn.Do("AUTH", redisPwd); err != nil {
				log.Println(err)
				return nil, err
			}
			return
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) (err error) {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err = c.Do("PING")
			return
		},
	}
}

func init() {
	pool = newRedisPool()
}

func GetPool() *redis.Pool {
	return pool
}
