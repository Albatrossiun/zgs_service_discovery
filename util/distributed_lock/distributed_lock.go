package distributed_lock

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
)

type DistributedLockManager struct {
	pool       *redis.Pool
	lockName   string
	expireTime int64
}

type DistributedLock struct {
	pool       *redis.Pool
	lockName   string
	lockValue  string
	expireTime int64
	c          chan int
	startTime  time.Time
}

func NewDistributedLockManager(pool *redis.Pool, lockName string, expireTime int64) *DistributedLockManager {
	return &DistributedLockManager{
		pool:       pool,
		lockName:   lockName,
		expireTime: expireTime,
	}
}

func (d *DistributedLockManager) GetLock() *DistributedLock {
	return &DistributedLock{
		pool:       d.pool,
		lockName:   d.lockName,
		lockValue:  uuid.NewString(),
		expireTime: d.expireTime,
		c:          make(chan int, 1),
		startTime:  time.Now(),
	}
}

func (d *DistributedLock) watchDog() {
	for {
		conn := d.pool.Get()

		select {
		case <-d.c:
			fmt.Print("看门狗退出")
			conn.Close()
			break
		default:
			passTime := time.Now().Sub(d.startTime).Seconds()
			if passTime > float64(d.expireTime)/3 {
				// 续期用的lua脚本
				lua := "local v=redis.call('keys', KEYS[1]); if v == nil then return 1 end v = redis.call('get', KEYS[1]); if v ~= KEYS[2] then return 1 end redis.call('EXPIRE', KEYS[1], KEYS[3]) return 0"
				code, err := conn.Do("eval", lua, 3, d.lockName, d.lockValue, d.expireTime)
				if code.(int) != 0 || err != nil {
					// 续期失败
					conn.Close()
					break
				}
			}
		}

		time.Sleep(1000 * time.Millisecond)
		conn.Close()
	}
}

func (d *DistributedLock) Lock() {
	conn := d.pool.Get()
	defer conn.Close()
	// 加锁
	// setnx 如果key不存在 才能设置成功
	_, err := conn.Do("setnx", d.lockName, d.lockValue, "ex", d.expireTime)
	// 如果加锁失败 在一个死循环里尝试加锁 直到加锁成功
	if err != nil {
		for {
			_, err = conn.Do("setnx", d.lockName, d.lockValue, "ex", d.expireTime)
			if err == nil {
				break
			}
			// 释放一下连接 避免一直占用连接耗尽redis连接池
			conn.Close()
			// 等待100毫秒
			time.Sleep(100 * time.Millisecond)
			conn = d.pool.Get()
		}
	}

	// 启动看门狗 持续看门
	go d.watchDog()
}

func (d *DistributedLock) TryLock() bool {
	conn := d.pool.Get()
	defer conn.Close()
	// 加锁
	// setnx 如果key不存在 才能设置成功
	_, err := conn.Do("setnx", d.lockName, d.lockValue, "ex", d.expireTime)
	// 如果加锁失败 则直接返回失败
	if err != nil {
		return false
	}

	// 启动看门狗 持续看门
	go d.watchDog()
	return true
}

func (d *DistributedLock) Unlock() {
	// 结束看门狗
	d.c <- 0
	// 释放锁
	conn := d.pool.Get()
	defer conn.Close()
	lua := "local v=redis.call('keys', KEYS[1]); if v == nil then return 1 end v = redis.call('get', KEYS[1]); if v ~= KEYS[2] then return 1 end redis.call('del', KEYS[1]) return 0"
	_, err := conn.Do("eval", lua, 2, d.lockName, d.lockValue)
	if err != nil {
		fmt.Println("解锁失败 ", err)
	}
}
