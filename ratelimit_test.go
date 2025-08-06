// Copyright 2024 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ratelimit

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"
)

func TestRateLimiter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("DefaultRateLimiter", func(t *testing.T) {
		r := gin.New()
		r.Use(New(Options{
			Rate:  rate.Every(time.Millisecond * 10),
			Burst: 1,
		}))
		r.GET("/", func(c *gin.Context) {
			c.String(http.StatusOK, "OK")
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusTooManyRequests, w.Code)
	})

	t.Run("CustomKeyFunc", func(t *testing.T) {
		r := gin.New()
		r.Use(New(Options{
			Rate:  rate.Every(time.Millisecond * 10),
			Burst: 1,
			KeyFunc: func(c *gin.Context) string {
				return c.Request.Header.Get("X-API-KEY")
			},
		}))
		r.GET("/", func(c *gin.Context) {
			c.String(http.StatusOK, "OK")
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Set("X-API-KEY", "test-key")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/", nil)
		req.Header.Set("X-API-KEY", "test-key")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusTooManyRequests, w.Code)
	})

	t.Run("CustomOnLimitExceeded", func(t *testing.T) {
		r := gin.New()
		r.Use(New(Options{
			Rate:  rate.Every(time.Millisecond * 10),
			Burst: 1,
			OnLimitExceeded: func(c *gin.Context, l *rate.Limiter) {
				c.String(http.StatusTeapot, "I'm a teapot")
			},
		}))
		r.GET("/", func(c *gin.Context) {
			c.String(http.StatusOK, "OK")
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusTeapot, w.Code)
		assert.Equal(t, "I'm a teapot", w.Body.String())
	})
}
