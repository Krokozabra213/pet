package httpPlatform

import (
	"time"

	"github.com/gin-contrib/cors"
)

// func corsMiddleware(c *gin.Context) {
// 	c.Header("Access-Control-Allow-Origin", "*")
// 	c.Header("Access-Control-Allow-Methods", "*")
// 	c.Header("Access-Control-Allow-Headers", "*")
// 	c.Header("Content-Type", "application/json")

// 	if c.Request.Method != "OPTIONS" {
// 		c.Next()
// 	} else {
// 		c.AbortWithStatus(http.StatusOK)
// 	}
// }

func getCorsConfig() cors.Config {
	return cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000", // frontend
		},
		AllowMethods: []string{
			"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"X-Requested-With",
		},
		ExposeHeaders: []string{
			"Content-Length",
		},
		AllowCredentials: true,
		MaxAge:           18 * time.Hour,
	}
}
