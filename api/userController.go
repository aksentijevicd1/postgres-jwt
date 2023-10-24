package api

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	db "github.com/aksentijevicd1/postgres-jwt/db"
	sqlcdb "github.com/aksentijevicd1/postgres-jwt/db/sqlc"
	"github.com/aksentijevicd1/postgres-jwt/helpers"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

var collection *pgxpool.Pool = db.InitDB()

type Production interface {
	sqlcdb.Querier
}

//store provides all functions to execute queries and transactions
type SQLProduction struct {
	*sqlcdb.Queries
	db *sql.DB
}

func NewProduction(db *sql.DB) Production {
	return &SQLProduction{
		db:      db,
		Queries: sqlcdb.New(db),
	}
}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprint("Email or password is incorrect!")
		check = false
	}

	return check, msg

}

func (server *Server) Signup() gin.HandlerFunc {

	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), time.Second*100)
		var user sqlcdb.User
		defer cancel()

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		conn, err := collection.Acquire(ctx)
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while acquiring a connection"})
			return
		}

		//defer conn.Release()
		var count = 0
		row := conn.QueryRow(ctx, "SELECT COUNT(*) FROM users WHERE email = $1", user.Email)
		if err = row.Scan(&count); err == sql.ErrNoRows {

			fmt.Println("There is no user with this particular email in database.")

		} else if err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred while checking for the email"})
			log.Panic(err)
			return

		} else {
			//uspesno
		}

		count = 0
		row = conn.QueryRow(ctx, "SELECT COUNT(*) FROM users WHERE phone = $1", user.Phone)
		if err = row.Scan(&count); err == sql.ErrNoRows {
			//no change. Meaning no user with that particular phone
		} else if err != nil {

			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred while checking for the phone"})
			return

		} else {
			// uspesno!
		}

		password := HashPassword(user.Password)
		user.Password = password

		user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		token, refreshToken, _ := helpers.GenerateAllTokens(user.Email, user.Firstname, user.Lastname, user.UserTypeID)
		user.Token = token
		user.RefreshToken = refreshToken

		_, errCreatingUser := server.production.CreateUser(ctx, sqlcdb.CreateUserParams{
			Firstname:    user.Firstname,
			Lastname:     user.Lastname,
			Password:     user.Password,
			Email:        user.Email,
			Phone:        user.Phone,
			Token:        user.Token,
			UserTypeID:   user.UserTypeID,
			RefreshToken: user.RefreshToken,
		})
		if errCreatingUser != nil {
			fmt.Println(errCreatingUser)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while creating user"})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, gin.H{"success": "You have successfully created user."})
	}

}

func (server *Server) Login() gin.HandlerFunc {
	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user sqlcdb.User
		var foundUser sqlcdb.User

		defer cancel()

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		conn, err := collection.Acquire(ctx)

		if err != nil {
			log.Panic("Error while trying to acquire connection from pool!")
		}
		defer conn.Release()
		row := conn.QueryRow(ctx, "SELECT * FROM users WHERE email = $1", user.Email)
		err = row.Scan(&foundUser.ID, &foundUser.Firstname,
			&foundUser.Lastname, &foundUser.Password, &foundUser.Email,
			&foundUser.Phone, &foundUser.Token, &foundUser.UserTypeID,
			&foundUser.RefreshToken, &foundUser.CreatedAt,
			&foundUser.UpdatedAt)
		if err == sql.ErrNoRows {

			c.JSON(http.StatusInternalServerError,
				gin.H{"error": "email is not found!"})
			log.Println("email is not found!")
			return

		} else if err != nil {

			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred while checking for the email"})
			return

		} else {

			// Email exists!

		}

		passwordIsValid, msg := VerifyPassword(user.Password, foundUser.Password)
		defer cancel()
		if passwordIsValid != true {
			c.JSON(http.StatusInternalServerError,
				gin.H{"error": msg})
			return
		}

		if foundUser.Email == "" {
			c.JSON(http.StatusInternalServerError,
				gin.H{"error": "user not found"})
			return
		}

		token, refreshToken, _ := helpers.GenerateAllTokens(foundUser.Email, foundUser.Firstname, foundUser.Lastname, foundUser.UserTypeID)
		helpers.UpdateAllTokens(token, refreshToken, foundUser.ID)
		row = conn.QueryRow(ctx, "SELECT * FROM users WHERE id = $1", foundUser.ID)
		if err = row.Scan(&foundUser.ID, &foundUser.Firstname,
			&foundUser.Lastname, &foundUser.Password, &foundUser.Email,
			&foundUser.Phone, &foundUser.Token, &foundUser.UserTypeID,
			&foundUser.RefreshToken, &foundUser.CreatedAt,
			&foundUser.UpdatedAt); err != nil {
			fmt.Println(user.ID)
			fmt.Println(user.ID)

			c.JSON(http.StatusInternalServerError, gin.H{"error": "user with this email is not found!"})
			log.Panic(err)
			return

		} else {

			c.JSON(http.StatusOK, foundUser)

		}

	}
}

func (server *Server) GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := helpers.CheckUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		recordPerPage, err := strconv.Atoi(c.DefaultQuery("recordPerPage", "10"))
		page, err1 := strconv.Atoi(c.DefaultQuery("page", "1"))
		startIndex, err2 := strconv.Atoi(c.DefaultQuery("startIndex", "0"))

		if err != nil || err1 != nil || err2 != nil || recordPerPage < 1 || page < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recordPerPage, page, or startIndex"})
			return
		}

		startIndex = (page - 1) * recordPerPage
		arg := sqlcdb.ListUsersParams{
			Limit:  int32(recordPerPage),
			Offset: int32(startIndex),
		}

		users, err := server.production.ListUsers(ctx, arg)
		defer cancel()

		if err != nil {
			log.Printf("Error querying users: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while listing user items"})
			return
		}
		c.JSON(http.StatusOK, users)
	}
}

func (server *Server) Getuser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userid, parseError := strconv.Atoi(c.Param("user_id"))
		if parseError != nil {
			log.Println("Error while parsing!")
			return
		}

		if err := helpers.CheckUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		user, err := server.production.GetUser(ctx, int64(userid))
		if err != nil {
			log.Printf("Error querying task: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while listing user"})
			return
		}

		c.JSON(http.StatusOK, user)

	}
}
