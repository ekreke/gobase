package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"sync"
	"testing"
	"time"
)

// getTestClient returns a Redis client for testing
func getTestClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   15, // Use separate DB for tests
	})
}

func TestLock_BasicLockUnlock(t *testing.T) {
	ctx := context.Background()
	client := getTestClient()
	defer client.Close()

	// Clean up before test
	client.Del(ctx, "lock:test-basic")

	lock := NewDistributedLock(client, "test-basic", time.Second*5).(*redisLock)

	// Test acquire lock
	acquired, err := lock.Lock(ctx)
	if err != nil {
		t.Fatalf("failed to acquire lock: %v", err)
	}
	if !acquired {
		t.Fatal("expected to acquire lock")
	}

	// Verify lock exists in Redis
	val, err := client.Get(ctx, "lock:test-basic").Result()
	if err != nil {
		t.Fatalf("failed to get lock from Redis: %v", err)
	}
	if val != lock.value {
		t.Fatalf("lock value mismatch: got %s, want %s", val, lock.value)
	}

	// Test unlock
	err = lock.Unlock(ctx)
	if err != nil {
		t.Fatalf("failed to unlock: %v", err)
	}

	// Verify lock is removed
	exists := client.Exists(ctx, "lock:test-basic").Val()
	if exists != 0 {
		t.Fatal("lock should be removed after unlock")
	}
}

func TestLock_DoubleLock(t *testing.T) {
	ctx := context.Background()
	client := getTestClient()
	defer client.Close()

	client.Del(ctx, "lock:test-double")

	lock := NewDistributedLock(client, "test-double", time.Second*5).(*redisLock)

	// First lock should succeed
	acquired, err := lock.Lock(ctx)
	if err != nil {
		t.Fatalf("first lock failed: %v", err)
	}
	if !acquired {
		t.Fatal("first lock should succeed")
	}

	// Second lock on same lock instance should succeed (idempotent)
	acquired, err = lock.Lock(ctx)
	if err != nil {
		t.Fatalf("second lock failed: %v", err)
	}
	if !acquired {
		t.Fatal("second lock should succeed (same instance)")
	}

	// Unlock
	err = lock.Unlock(ctx)
	if err != nil {
		t.Fatalf("unlock failed: %v", err)
	}
}

func TestLock_ConcurrentLock(t *testing.T) {
	ctx := context.Background()
	client := getTestClient()
	defer client.Close()

	client.Del(ctx, "lock:test-concurrent")

	lock1 := NewDistributedLock(client, "test-concurrent", time.Second*5)
	lock2 := NewDistributedLock(client, "test-concurrent", time.Second*5)

	var wg sync.WaitGroup
	wg.Add(2)

	results := make([]bool, 2)
	var errors [2]error

	// Both try to acquire the lock
	go func(idx int) {
		defer wg.Done()
		acquired, err := lock1.Lock(ctx)
		results[idx] = acquired
		errors[idx] = err
	}(0)

	go func(idx int) {
		defer wg.Done()
		acquired, err := lock2.Lock(ctx)
		results[idx] = acquired
		errors[idx] = err
	}(1)

	wg.Wait()

	if errors[0] != nil {
		t.Fatalf("goroutine 0 error: %v", errors[0])
	}
	if errors[1] != nil {
		t.Fatalf("goroutine 1 error: %v", errors[1])
	}

	// Only one should acquire the lock
	acquiredCount := 0
	if results[0] {
		acquiredCount++
	}
	if results[1] {
		acquiredCount++
	}

	if acquiredCount != 1 {
		t.Fatalf("expected exactly 1 lock acquired, got %d", acquiredCount)
	}

	t.Logf("goroutine 0 acquired: %v, goroutine 1 acquired: %v", results[0], results[1])
}

func TestLock_UnlockWithoutLock(t *testing.T) {
	ctx := context.Background()
	client := getTestClient()
	defer client.Close()

	client.Del(ctx, "lock:test-unlock-without")

	lock := NewDistributedLock(client, "test-unlock-without", time.Second*5)

	// Unlock without acquiring should fail
	err := lock.Unlock(ctx)
	if err == nil {
		t.Fatal("expected error when unlocking without lock")
	}

	expectedErrMsg := "lock not held by current session"
	if err.Error() != expectedErrMsg {
		t.Fatalf("unexpected error message: got %s, want %s", err.Error(), expectedErrMsg)
	}
}

func TestLock_LockExpiration(t *testing.T) {
	ctx := context.Background()
	client := getTestClient()
	defer client.Close()

	client.Del(ctx, "lock:test-expiration")

	shortExpiration := 500 * time.Millisecond
	lock := NewDistributedLock(client, "test-expiration", shortExpiration).(*redisLock)

	// Acquire lock
	acquired, err := lock.Lock(ctx)
	if err != nil {
		t.Fatalf("failed to acquire lock: %v", err)
	}
	if !acquired {
		t.Fatal("failed to acquire lock")
	}

	// Get TTL (Redis returns approximate value in seconds)
	ttl := client.TTL(ctx, "lock:test-expiration").Val()
	if ttl < 0 {
		t.Fatalf("unexpected TTL: got %v", ttl)
	}
	t.Logf("TTL: %v", ttl)

	// Wait for watchdog to renew
	time.Sleep(shortExpiration/2 + 100*time.Millisecond)

	// Check if lock still exists (watchdog should have renewed)
	exists := client.Exists(ctx, "lock:test-expiration").Val()
	if exists != 1 {
		t.Fatal("lock should still exist after watchdog renewal")
	}

	// Unlock to cleanup
	err = lock.Unlock(ctx)
	if err != nil {
		t.Fatalf("unlock failed: %v", err)
	}
}

func TestLock_MultipleLocksDifferentKeys(t *testing.T) {
	ctx := context.Background()
	client := getTestClient()
	defer client.Close()

	client.Del(ctx, "lock:test-key1")
	client.Del(ctx, "lock:test-key2")

	lock1 := NewDistributedLock(client, "test-key1", time.Second*5)
	lock2 := NewDistributedLock(client, "test-key2", time.Second*5)

	// Both should succeed since they have different keys
	acquired1, err := lock1.Lock(ctx)
	if err != nil {
		t.Fatalf("lock1 failed: %v", err)
	}
	if !acquired1 {
		t.Fatal("lock1 should succeed")
	}

	acquired2, err := lock2.Lock(ctx)
	if err != nil {
		t.Fatalf("lock2 failed: %v", err)
	}
	if !acquired2 {
		t.Fatal("lock2 should succeed")
	}

	// Clean up
	lock1.Unlock(ctx)
	lock2.Unlock(ctx)

	// Verify both are removed
	exists := client.Exists(ctx, "lock:test-key1", "lock:test-key2").Val()
	if exists != 0 {
		t.Fatal("both locks should be removed")
	}
}

func TestLock_ValueMismatch(t *testing.T) {
	ctx := context.Background()
	client := getTestClient()
	defer client.Close()

	client.Del(ctx, "lock:test-mismatch")

	// Manually set a lock with a different value
	client.Set(ctx, "lock:test-mismatch", "someone-elses-value", time.Second*10)

	lock := NewDistributedLock(client, "test-mismatch", time.Second*5).(*redisLock)
	lock.value = "my-different-value"

	// Try to unlock - should fail due to value mismatch
	// First, we need to set held=true to pass the initial check
	lock.held = true
	lock.watchdogStop = make(chan struct{})
	lock.watchdogDone = make(chan struct{})
	close(lock.watchdogStop)
	close(lock.watchdogDone)

	err := lock.Unlock(ctx)
	if err == nil {
		t.Fatal("expected error when unlocking with wrong value")
	}

	if err.Error() != "lock value mismatch - possibly expired or stolen" {
		t.Logf("got error: %v", err)
	}

	// The original lock should still exist in Redis
	exists := client.Exists(ctx, "lock:test-mismatch").Val()
	if exists != 1 {
		t.Fatal("original lock should still exist")
	}

	// Cleanup
	client.Del(ctx, "lock:test-mismatch")
}

func TestLock_RetryAfterExpiry(t *testing.T) {
	ctx := context.Background()
	client := getTestClient()
	defer client.Close()

	client.Del(ctx, "lock:test-retry")

	lock := NewDistributedLock(client, "test-retry", 200*time.Millisecond).(*redisLock)

	// First lock
	acquired, err := lock.Lock(ctx)
	if err != nil {
		t.Fatalf("first lock failed: %v", err)
	}
	if !acquired {
		t.Fatal("first lock should succeed")
	}

	// Wait for lock to expire
	time.Sleep(300 * time.Millisecond)

	// Try to lock again with the same instance
	// The old lock should have expired in Redis
	// but lock.held is still true locally, so it will return true
	acquired, err = lock.Lock(ctx)
	if err != nil {
		t.Fatalf("second lock failed: %v", err)
	}
	if !acquired {
		t.Fatal("second lock should succeed (held=true locally)")
	}

	// Clean up - unlock
	err = lock.Unlock(ctx)
	if err != nil {
		t.Fatalf("unlock failed: %v", err)
	}
}

func TestLock_WatchdogRenewal(t *testing.T) {
	ctx := context.Background()
	client := getTestClient()
	defer client.Close()

	client.Del(ctx, "lock:test-watchdog")

	expiration := 5 * time.Second
	lock := NewDistributedLock(client, "test-watchdog", expiration).(*redisLock)

	// Acquire lock
	acquired, err := lock.Lock(ctx)
	if err != nil {
		t.Fatalf("failed to acquire lock: %v", err)
	}
	if !acquired {
		t.Fatal("failed to acquire lock")
	}

	// Get initial TTL
	ttl1 := client.TTL(ctx, "lock:test-watchdog").Val()
	t.Logf("Initial TTL: %v", ttl1)

	// Wait for watchdog to renew at least once
	// Watchdog runs every expiration/2 = 2.5s
	sleepTime := expiration/2 + 500*time.Millisecond
	t.Logf("Sleeping %v for watchdog to renew...", sleepTime)
	time.Sleep(sleepTime)

	// Check TTL again
	ttl2 := client.TTL(ctx, "lock:test-watchdog").Val()
	t.Logf("TTL after watchdog: %v", ttl2)

	// TTL should still be high (lock was renewed)
	// If no renewal happened, TTL would be around 5s - 3s = 2s
	// If renewal happened, TTL should be back around 5s
	if ttl2 < 3*time.Second {
		t.Fatalf("watchdog renewal failed: TTL after %v sleep is only %v, want > 3s", sleepTime, ttl2)
	}

	// Clean up
	err = lock.Unlock(ctx)
	if err != nil {
		t.Fatalf("unlock failed: %v", err)
	}

	// Verify lock is removed
	exists := client.Exists(ctx, "lock:test-watchdog").Val()
	if exists != 0 {
		t.Fatal("lock should be removed after unlock")
	}
}

// Benchmark test
func BenchmarkLock_LockUnlock(b *testing.B) {
	ctx := context.Background()
	client := getTestClient()
	defer client.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("bench-lock-%d", i%100) // Use 100 different keys
		client.Del(ctx, "lock:"+key)

		lock := NewDistributedLock(client, key, time.Second*5)
		lock.Lock(ctx)
		lock.Unlock(ctx)
	}
}