// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0

package db

import (
	"time"
)

type Category struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type Status struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type Task struct {
	ID           int64     `json:"id"`
	Title        string    `json:"title"`
	Illustration string    `json:"illustration"`
	CategoryID   int64     `json:"category_id"`
	StatusID     int64     `json:"status_id"`
	CreatedAt    time.Time `json:"created_at"`
}

type User struct {
	ID           int64     `json:"id"`
	Firstname    string    `json:"firstname"`
	Lastname     string    `json:"lastname"`
	Password     string    `json:"password"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	Token        string    `json:"token"`
	UserTypeID   int64     `json:"user_type_id"`
	RefreshToken string    `json:"refresh_token"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Usertype struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}