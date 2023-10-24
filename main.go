package main

import (
	"database/sql"
	"log"

	"github.com/aksentijevicd1/postgres-jwt/api"
	db "github.com/aksentijevicd1/postgres-jwt/db/sqlc"
	"github.com/aksentijevicd1/postgres-jwt/util"
	_ "github.com/lib/pq"
)

var testQueries *db.Queries

func main() {
	/*router := gin.New()
	router.Use(gin.Logger())

	routes.AuthRoutes(router)
	routes.UserRoutes(router)

	router.GET("/api-1", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "Access granted for api-1"})
	})

	router.GET("/api-2", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "Access granted for api-2"})
	})

	router.Run(":8080")
	*/
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("can't load config", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("can't connect to db", err)
	}

	production := db.NewProduction(conn)
	server := api.NewServer(production)
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server: ", err)

	}
}
