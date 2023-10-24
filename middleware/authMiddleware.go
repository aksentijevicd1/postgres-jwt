package middleware

import (
	"fmt"
	"net/http"

	"github.com/aksentijevicd1/postgres-jwt/helpers"
	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("No authorization header provided")})
			c.Abort()
			return
		}

		claims, err := helpers.ValidateToken(clientToken)
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("first_name", claims.Firstname)
		c.Set("last_name", claims.Lastname)
		c.Set("user_type_id", claims.UserTypeID)
		c.Next()
	}
}
