package models

import "time"

type User struct {
	FirstName string `json:"first_name" binding:"required,min=2,max=50"`
	Email     string `json:"email" binding:"required,email"`
}

type VerifyEmail struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type GetUser struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}
