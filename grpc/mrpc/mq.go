package mrpc

import (
	"context"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"time"
)

//简易的redis pub/sub

func init() {
	InitRedis()
}

var redisClient *redis.Pool

const healthCheckPeriod = time.Second

type CallbackMessage interface {
	HandleMessage(string, []byte) ([]byte, error)
}

type SubManager struct {
	Callback CallbackMessage
}

var SubManagerInstance *SubManager

func InitSubManager(message CallbackMessage) error {
	if SubManagerInstance != nil {
		return fmt.Errorf("init error ")
	}

	SubManagerInstance = &SubManager{
		Callback: message,
	}
	return nil
}

func InitRedis() {
	redisClient = &redis.Pool{
		// 最大空闲链接
		MaxIdle: 10,
		// 最大激活链接
		MaxActive: 10,
		// 最大空闲链接等待时间
		IdleTimeout: 5 * time.Second,
		Dial: func() (redis.Conn, error) {
			rc, err := redis.Dial("tcp", "127.0.0.1:6379")
			if err != nil {
				return nil, err
			}
			rc.Do("SELECT", 0)
			fmt.Println("USE DB", 0)

			return rc, nil
		},
	}
}

//订阅
func Subscribe(channel string) (context.CancelFunc, error) {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		defer func() {
			fmt.Println("Subscribe go func close")
		}()

		for {
			minticker := time.After(5 * time.Second)

			psc := redis.PubSubConn{Conn: redisClient.Get()}
			if err := psc.Subscribe(channel); err != nil { //订阅此频道
				select {
				case <-minticker:
				}
				continue
				fmt.Errorf("error")
			}

			// Start a goroutine to receive notifications from the server.
			done := make(chan error, 1)

			go func() {
				for {
					switch n := psc.Receive().(type) { //接收
					case error:
						done <- n
						return
					case redis.Message:
						if _, err := SubManagerInstance.Callback.HandleMessage(n.Channel, n.Data); err != nil {
							done <- err
							return
						}
					case redis.Subscription:
						if n.Count == 0 {
							// all channels are unsubscribed
							done <- nil
							return
						}
					}
				}
			}()

			//心跳检测
			ticker := time.NewTicker(healthCheckPeriod)

		loop:
			for {
				select {
				case <-ticker.C:
					// Send ping to test health of connection and server. If
					// corresponding pong is not received, then receive on the
					// connection will timeout and the receive goroutine will exit.
					if err := psc.Ping(""); err != nil {
						break loop
					}
				case <-ctx.Done():
					break loop
				case <-done:
					break loop
					// Return error from the receive goroutine.
					//return nil,err
				}
			}
		}
	}()

	return cancel, nil
}

//发布
func Publish(channel, message string) (int, error) {
	conn := redisClient.Get()
	defer conn.Close()

	n, err := redis.Int(conn.Do("PUBLISH", channel, message))
	if err != nil {
		return 0, fmt.Errorf("redis publish %s %s, err: %v", channel, message, err)
	}

	return n, nil
}
