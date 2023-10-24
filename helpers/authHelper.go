package helpers

import (
	"context"
	"errors"
	"log"
	"time"

	database "github.com/aksentijevicd1/postgres-jwt/db"
	"github.com/gin-gonic/gin"
)

func CheckUserType(c *gin.Context, role string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*90)
	defer cancel()
	userTypeID := c.GetInt64("user_type_id")
	err = nil
	conn, DBErr := database.GetDB().Acquire(ctx)
	if DBErr != nil {
		log.Fatal("Error while trying to acquire conn from pool!")
		return
	}
	query := `SELECT name FROM usertypes WHERE id = (SELECT id FROM users WHERE id = $1)`
	row := conn.QueryRow(ctx, query, userTypeID)
	var userTypeString string
	err1 := row.Scan(&userTypeString)

	if err1 != nil {
		log.Fatal("There is no particular role in our DB")
		return
	}

	if userTypeString != role {
		err = errors.New("Unauthorized access to this resource!")
		return err
	}

	return err

}
