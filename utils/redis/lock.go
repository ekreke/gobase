package redis

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// Lock defines the distributed lock interface
type Lock interface {
	Lock(ctx context.Context) (bool, error)
	Unlock(ctx context.Context) error
}

// redisLock implements the Lock interface
type redisLock struct {
	rdb        redis.Cmdable
	key        string // resource name
	value      string // uuid
	expiration time.Duration

	mu           sync.Mutex
	held         bool          // if the lock is held
	watchdogStop chan struct{} // notify the watchdog stop renew
	watchdogDone chan struct{} // confirm the watchdog stopped
}

// NewDistributedLock creates a new Redis distributed lock
func NewDistributedLock(rdb redis.Cmdable, key string, expiration time.Duration) Lock {
	return &redisLock{
		rdb:        rdb,
		key:        "lock:" + key,
		expiration: expiration,
	}
}

// Lock try to get lock , if failed instantly return
func (l *redisLock) Lock(ctx context.Context) (bool, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.held {
		return true, nil
	}

	l.value = uuid.New().String()

	// use set nx ex to get redis lock
	ok, err := l.rdb.SetNX(ctx, l.key, l.value, l.expiration).Result()
	if err != nil {
		return false, err
	}

	if !ok {
		return false, nil // lock is held by others
	}

	// if successfully get the lock , start the watchdog goroutine
	l.held = true
	l.watchdogStop = make(chan struct{})
	l.watchdogDone = make(chan struct{})
	go l.watchdog()

	return true, nil
}

// watchdog to renew goroutine automatically
func (l *redisLock) watchdog() {
	defer close(l.watchdogDone)

	interval := l.expiration / 2
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-l.watchdogStop:
			return
		case <-ticker.C:
			l.renew()
		}
	}
}

// Unlock try to unlock the redis lock
func (l *redisLock) Unlock(ctx context.Context) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if !l.held {
		return errors.New("lock not held by current session")
	}

	// Safely close the watchdog stop channel
	select {
	case <-l.watchdogStop:
		// Already closed
	default:
		close(l.watchdogStop)
	}
	<-l.watchdogDone

	script := redis.NewScript(`
		if redis.call("GET", KEYS[1]) == ARGV[1] then
			return redis.call("DEL", KEYS[1])
		else
			return 0
		end
	`)

	res, err := script.Run(ctx, l.rdb, []string{l.key}, l.value).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}

	if res != nil && res.(int64) == 0 {
		return errors.New("lock value mismatch - possibly expired or stolen")
	}

	l.held = false

	return nil
}

// renew try to renew the lock
func (l *redisLock) renew() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	script := redis.NewScript(`
		if redis.call("GET", KEYS[1]) == ARGV[1] then
			return redis.call("PEXPIRE", KEYS[1], ARGV[2])
		else
			return 0
		end
	`)
	renewMs := int64(l.expiration / time.Millisecond)
	_, _ = script.Run(ctx, l.rdb, []string{l.key}, l.value, renewMs).Result()
}
