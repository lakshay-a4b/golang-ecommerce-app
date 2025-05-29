package models

import (
	"time"
)

type User struct {
	UserId    string    `json:"userId"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
}

type UserEvent struct {
	Timestamp string `json:"date"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Action    string `json:"action"`
	UserId    string `json:"userId"`
}
