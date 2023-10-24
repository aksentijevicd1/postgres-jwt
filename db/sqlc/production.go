package db

import "database/sql"

type Production interface {
	Querier
}

//store provides all functions to execute queries and transactions
type SQLProduction struct {
	*Queries
	db *sql.DB
}

func NewProduction(db *sql.DB) Production {
	return &SQLProduction{
		db:      db,
		Queries: New(db),
	}
}
