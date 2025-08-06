// Copyright 2024 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Package ratelimit provides a rate limiting middleware for the Gin framework.
// It is based on the token bucket algorithm and can be used to limit the number
// of requests a client can make in a given amount of time.
package ratelimit

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// Options contains the configuration for the rate limiting middleware.
type Options struct {
	// Rate is the token generation rate (tokens per second).
	// It determines how many requests are allowed per second.
	Rate rate.Limit

	// Burst is the bucket size (maximum burst of requests).
	// It determines the maximum number of requests that can be
	// handled in a short burst.
	Burst int

	// KeyFunc is a function to generate a key for rate limiting.
	// The key is used to identify a client and apply the rate limit
	// to that client. If nil, the client's IP address is used.
	KeyFunc func(*gin.Context) string

	// Store is the storage for rate limiters.
	// It is used to store the rate limiters for each client.
	// If nil, a default in-memory store is used.
	Store Store

	// OnLimitExceeded is a handler called when the rate limit is exceeded.
	// It can be used to customize the response sent to the client when
	// the rate limit is exceeded. If nil, a default handler that sends a
	// 429 Too Many Requests response is used.
	OnLimitExceeded func(*gin.Context, *rate.Limiter)
}

// Store is the interface for storing rate limiters.
// It can be implemented to use different storage backends,
// such as in-memory, Redis, or others.
type Store interface {
	// Get retrieves a rate limiter from the store for the given key.
	Get(key string) (*rate.Limiter, bool)
	// Set adds a rate limiter to the store for the given key.
	Set(key string, limiter *rate.Limiter)
}

// New creates a new rate limiting middleware with the given options.
func New(opts Options) gin.HandlerFunc {
	// Set default options if not provided.
	if opts.KeyFunc == nil {
		opts.KeyFunc = func(c *gin.Context) string {
			return c.ClientIP()
		}
	}
	if opts.Store == nil {
		opts.Store = newMemoryStore()
	}
	if opts.OnLimitExceeded == nil {
		opts.OnLimitExceeded = func(c *gin.Context, l *rate.Limiter) {
			c.String(http.StatusTooManyRequests, "Too Many Requests")
		}
	}

	return func(c *gin.Context) {
		// Generate a key for the client.
		key := opts.KeyFunc(c)
		// Get the rate limiter for the client from the store.
		limiter, exists := opts.Store.Get(key)
		if !exists {
			// If the rate limiter does not exist, create a new one
			// and add it to the store.
			limiter = rate.NewLimiter(opts.Rate, opts.Burst)
			opts.Store.Set(key, limiter)
		}

		// Check if the client has exceeded the rate limit.
		if !limiter.Allow() {
			// If the rate limit is exceeded, call the OnLimitExceeded handler.
			opts.OnLimitExceeded(c, limiter)
			c.Abort()
			return
		}

		// If the rate limit is not exceeded, continue to the next handler.
		c.Next()
	}
}

// memoryStore is an in-memory implementation of the Store interface.
// It uses a map to store the rate limiters for each client.
type memoryStore struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
}

// newMemoryStore creates a new in-memory store.
func newMemoryStore() *memoryStore {
	return &memoryStore{
		limiters: make(map[string]*rate.Limiter),
	}
}

// Get retrieves a rate limiter from the store.
func (s *memoryStore) Get(key string) (*rate.Limiter, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	limiter, exists := s.limiters[key]
	return limiter, exists
}

// Set adds a rate limiter to the store.
func (s *memoryStore) Set(key string, limiter *rate.Limiter) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.limiters[key] = limiter
}
