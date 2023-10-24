package helpers

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	database "github.com/aksentijevicd1/postgres-jwt/db"
	"github.com/dgrijalva/jwt-go"
	"github.com/jackc/pgx/v4/pgxpool"
)

type SignedDetails struct {
	Email      string
	Firstname  string
	Lastname   string
	UserTypeID int64
	jwt.StandardClaims
}

var Collection *pgxpool.Pool = database.GetDB()
var SECRET_KEY string = os.Getenv("SECRET_KEY")

func GenerateAllTokens(email string, firstName string, lastName string, userType int64) (signedToken string, signedRefreshToken string, err error) {
	claims := &SignedDetails{
		Email:      email,
		Firstname:  firstName,
		Lastname:   lastName,
		UserTypeID: userType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Panic(err)
		return
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Panic(err)
		return
	}

	return token, refreshToken, err
}

func UpdateAllTokens(signedToken string, signedRefreshToken string, ID int64) {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	defer cancel()

	query := `UPDATE users SET token = $1, refresh_token = $2, updated_at = $3 WHERE id = $4`

	conn, err := database.GetDB().Acquire(ctx)
	defer cancel()

	if err != nil {
		log.Panic(err)
		return
	}

	_, err = conn.Exec(ctx, query, signedToken, signedRefreshToken, time.Now(), ID)

	if err != nil {

		log.Panic(err)
		return

	}

	return

}

func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)

	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = fmt.Sprintf("the token is invalid")
		msg = err.Error()
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = fmt.Sprintf("token is expired")
		msg = err.Error()
		return
	}
	return claims, msg
}
