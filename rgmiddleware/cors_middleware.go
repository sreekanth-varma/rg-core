package rgmiddleware

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	isEnabled := os.Getenv("middleware_cors_enabled")
	log.Println("cors enables ? ", os.Getenv("middleware_cors_enabled"))
	log.Println("cors enabled")
	return func(c *gin.Context) {
		if isEnabled != "Y" {
			log.Println("Cors is in disabled")
			c.Next()
			return
		}

		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Encrypted")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "OPTIONS,GET,POST,PUT,PATCH,DELETE")
		log.Println("cors .................")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
