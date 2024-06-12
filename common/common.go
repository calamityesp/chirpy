package common

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

type User struct {
	Id                 int    `json:"id"`
	Email              string `json:"email"`
	Password           string `json:"password"`
	Token              string `json:"jwtToken1"`
	Expires_in_seconds int    `json:"expires_in_seconds"`
}
