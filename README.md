# Gin Rate Limiting Middleware

This document provides instructions on how to use the rate limiting middleware for the Gin framework.

## Installation

To use the rate limiting middleware, you first need to import it into your project:

```go
import "github.com/gin-gonic/gin/middleware/ratelimit"
```

## Usage

The rate limiting middleware can be used as a global middleware, for a specific route group, or for a single route.

### Global Rate Limiting

To apply rate limiting to all routes, use the `New` function as a global middleware:

```go
import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/middleware/ratelimit"
	"golang.org/x/time/rate"
)

func main() {
	r := gin.Default()

	// Apply rate limiting to all requests
	r.Use(ratelimit.New(ratelimit.Options{
		Rate:  rate.Every(time.Second),
		Burst: 10,
	}))

	// ... your routes
}
```

### Route Group Rate Limiting

To apply rate limiting to a specific group of routes, use the `New` function within a route group:

```go
func main() {
	r := gin.Default()

	apiGroup := r.Group("/api")
	apiGroup.Use(ratelimit.New(ratelimit.Options{
		Rate:  rate.Every(time.Minute),
		Burst: 100,
	}))

	// ... routes within the API group
}
```

### Customizing Rate Limiting

The `Options` struct allows you to customize the rate limiting behavior:

- `Rate`: The rate at which tokens are generated (e.g., `rate.Every(time.Second)` for one token per second).
- `Burst`: The maximum number of tokens that can be stored in the bucket.
- `KeyFunc`: A function to generate a unique key for each client. By default, the client's IP address is used.
- `Store`: The storage backend for rate limiters. By default, an in-memory store is used. You can also use a Redis-based store for distributed rate limiting.
- `OnLimitExceeded`: A function that is called when a client exceeds the rate limit. By default, a `429 Too Many Requests` response is sent.

### Using a Redis Store

To use a Redis-based store for distributed rate limiting, you need to create a `redis.Client` and pass it to the `NewRedisStore` function:

```go
import (
	"github.com/go-redis/redis/v8"
	"github.com/gin-gonic/gin/middleware/ratelimit"
)

func main() {
	// ...

	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	r.Use(ratelimit.New(ratelimit.Options{
		// ...
		Store: ratelimit.NewRedisStore(redisClient),
	}))

	// ...
}
```