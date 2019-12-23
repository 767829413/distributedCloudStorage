package conn

import (
	. "distributedCloudStorage/common"
	"github.com/garyburd/redigo/redis"
	"log"
	"time"
)

var (
	pool *redis.Pool
)

func newRedisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     50,
		MaxActive:   30,
		IdleTimeout: 300 * time.Second,
		Dial: func() (conn redis.Conn, err error) {
			if conn, err = redis.Dial("tcp", CacheHost); err != nil {
				log.Println(err.Error())
				return
			}
			if _, err = conn.Do("AUTH", CachePwd); err != nil {
				log.Println(err.Error())
				return
			}
			return
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) (err error) {
			if time.Since(t) < time.Minute {
				return
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
