package api

import (
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine, server *Server) {
	incomingRoutes.POST("users/signup", server.Signup())
	incomingRoutes.POST("users/login", server.Login())
}
