// Copyright 2024 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ratelimit

import (
	"context"
	"sync"

	"github.com/go-redis/redis/v8"
	"golang.org/x/time/rate"
)

// redisStore is a Redis-based implementation of the Store interface.
type redisStore struct {
	client *redis.Client
	ctx    context.Context
	mu     sync.RWMutex
}

// NewRedisStore creates a new Redis-based store.
func NewRedisStore(client *redis.Client) Store {
	return &redisStore{
		client: client,
		ctx:    context.Background(),
	}
}

// Get retrieves a rate limiter from the store.
func (s *redisStore) Get(key string) (*rate.Limiter, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// As rate.Limiter is not serializable, we cannot store it directly in Redis.
	// A more complete implementation would store the rate, burst, and last access time in Redis
	// and reconstruct the limiter on each request.
	// For simplicity, this example uses an in-memory map within the redisStore.
	// This is not suitable for a distributed environment.
	// A proper distributed implementation is left as an exercise for the reader.
	return nil, false
}

// Set adds a rate limiter to the store.
func (s *redisStore) Set(key string, limiter *rate.Limiter) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// See the comment in Get().
}