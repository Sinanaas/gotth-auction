package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	limitergin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

func RateLimiter(rate limiter.Rate) gin.HandlerFunc {
	store := memory.NewStore()
	instance := limiter.New(store, rate)
	return limitergin.NewMiddleware(instance)
}
