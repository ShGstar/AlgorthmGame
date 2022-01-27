package redis_mq

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"log"
)

type RedisProcessor struct {
	mServerAddress string
	mPassword      string
	mDBNum         int
	mThreadsNum    int //开启N个conn N个协程并发
	connMap        map[*redis.Conn]bool
	Done           chan interface{}
	close          chan interface{}
	taksFIFO       chan *ConnTask
}

type ConnTask struct {
	cmd           string
	args          []interface{}
	TaskResult    chan interface{}
	luaTaskScript *redis.Script
}

func (m *RedisProcessor) RedisDial() {
	for i := 0; i < m.mThreadsNum; i++ {
		conn, err := redis.Dial("tcp", m.mServerAddress,
			redis.DialPassword(m.mPassword), redis.DialDatabase(m.mDBNum))
		if err != nil {
			log.Panicf("Redis Connect Error : %v", err)
		}
		m.connMap[&conn] = true

		go m.ConnRecover(&conn, i)
	}
	fmt.Printf("RedisDial %d success \n", m.mThreadsNum)
}

func (m *RedisProcessor) ConnRecover(conn *redis.Conn, goruntinueNum int) {
	for {
		select {
		case task := <-m.taksFIFO:
			if task.luaTaskScript != nil {
				value, err := (*task.luaTaskScript).Do(*conn, task.args...)
				if err != nil {
					fmt.Println("lua conn do err :", err)
				}
				task.TaskResult <- value
				fmt.Println("lua conn goruntinueNum :", goruntinueNum, " is processing ,", task.cmd, " ", task.args)
			} else {
				value, err := (*conn).Do(task.cmd, task.args...)
				if err != nil {
					fmt.Println("conn do err :", err)
				}
				task.TaskResult <- value
				fmt.Println("goruntinueNum :", goruntinueNum, " is processing ,", task.cmd, " ", task.args)
			}
		case <-m.close:
			m.Done <- "close conn Goroutinue"
			fmt.Println("goruntinueNum is close :", goruntinueNum)
			return
		}
	}
}

func (m *RedisProcessor) PushTaskDo(cmd string, args ...interface{}) *ConnTask {
	task := &ConnTask{}
	task.cmd = cmd
	task.args = args
	task.TaskResult = make(chan interface{}, 1)
	m.taksFIFO <- task
	return task
}

func (m *RedisProcessor) Close() {
	close(m.close)
	closed := 0
	for conn := range m.connMap {
		err := (*conn).Close()
		if err != nil {
			log.Panic(err)
			continue
		}
		closed++
	}

	//wait conn Goroutinue done
	for i := 0; i < m.mThreadsNum; i++ {
		<-m.Done
	}
	fmt.Printf("Redis connMap %d close successful \n", closed)
}
