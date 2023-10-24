package db

import (
	"context"
	"log"

	"github.com/aksentijevicd1/postgres-jwt/util"
	"github.com/jackc/pgx/v4/pgxpool"
)

var db *pgxpool.Pool

func InitDB() *pgxpool.Pool {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("can't load config", err)
	}

	// Create a database connection pool
	pool, err := pgxpool.Connect(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("can't connect to db", err)
	}

	db = pool

	return db
}

func GetDB() *pgxpool.Pool {
	return db
}
