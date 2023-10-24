package api

import (
	db "github.com/aksentijevicd1/postgres-jwt/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	production db.Production
	router     *gin.Engine
}

func NewServer(production db.Production) *Server {
	server := &Server{production: production}
	router := gin.Default()
	UserRoutes(router, server)
	AuthRoutes(router, server)
	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
