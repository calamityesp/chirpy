package common

import (
	"time"
)

type Chirp struct {
	ID        int    `json:"id"`
	Body      string `json:"body"`
	Author_Id int    `json:"author_id"`
}

type User struct {
	Id                        int    `json:"id"`
	Email                     string `json:"email"`
	Password                  string `json:"password"`
	Token                     string `json:"token"`
	RefreshToken              string `json:"refresh_token"`
	Expires_in_seconds        int    `json:"expires_in_seconds"`
	Refresh_token_expire_time time.Time
	Chirpy_Red                bool `json:"is_chirpy_red"`
}

type Polka struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}
