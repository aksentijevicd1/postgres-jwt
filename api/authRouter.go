package api

import (
	"github.com/aksentijevicd1/postgres-jwt/middleware"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes *gin.Engine, server *Server) {
	incomingRoutes.Use(middleware.Authenticate())

	incomingRoutes.GET("/users", server.GetUsers())
	incomingRoutes.GET("/users/:user_id", server.Getuser())
}
