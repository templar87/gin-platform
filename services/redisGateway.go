package services

import (
	"flag"
	"time"

	"github.com/garyburd/redigo/redis"
)

var (
	RedisPool     *redis.Pool
	redisServer   = flag.String("redisServer", "127.0.0.1:6379", "")
	redisPassword = flag.String("redisPassword", "XXXXXXXX", "")
)

func InitRedis() {
	flag.Parse()
	RedisPool = newRedisPool(*redisServer, *redisPassword)
}
func newRedisPool(server, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     5,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			/*if _, err := c.Do("AUTH", password); err != nil {
			    c.Close()
			    return nil, err
			}*/
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func SendCmdToRedis(cmd, key string, value interface{}) {
	redisGateway := RedisPool.Get()
	defer redisGateway.Close()

	_, _ = redisGateway.Do(cmd, key, value)
}
