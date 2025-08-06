package main

import (
	"net/http"
	"time"

	"github.com/gin-contrib/ratelimit"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"golang.org/x/time/rate"
)

func main() {
	// Default usage
	app := gin.Default()
	app.Use(ratelimit.New(ratelimit.Options{
		Rate:  rate.Every(time.Second),
		Burst: 1,
	}))
	app.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})

	// Custom usage
	customApp := gin.Default()
	customApp.Use(ratelimit.New(ratelimit.Options{
		Rate:  rate.Every(time.Minute),
		Burst: 5,
		KeyFunc: func(c *gin.Context) string {
			return c.GetHeader("X-API-KEY")
		},
		OnLimitExceeded: func(c *gin.Context, l *rate.Limiter) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"message": "Too many requests",
			})
		},
	}))
	customApp.GET("/custom", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, Custom World!")
	})

	// Redis usage
	redisApp := gin.Default()
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	redisApp.Use(ratelimit.New(ratelimit.Options{
		Rate:  rate.Every(time.Second),
		Burst: 1,
		Store: ratelimit.NewRedisStore(redisClient),
	}))
	redisApp.GET("/redis", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, Redis World!")
	})

	// Per-route usage
	perRouteApp := gin.Default()
	rateLimiter := ratelimit.New(ratelimit.Options{
		Rate:  rate.Every(time.Second),
		Burst: 1,
	})
	perRouteApp.GET("/limited", rateLimiter, func(c *gin.Context) {
		c.String(http.StatusOK, "This is a limited route")
	})
	perRouteApp.GET("/unlimited", func(c *gin.Context) {
		c.String(http.StatusOK, "This is an unlimited route")
	})

	go func() {
		app.Run(":8080")
	}()
	go func() {
		customApp.Run(":8081")
	}()
	go func() {
		redisApp.Run(":8082")
	}()
	perRouteApp.Run(":8083")
}
