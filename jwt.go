package main

import (
	"log"
	"strconv"
	"time"

	"github.com/calamityesp/chirpy/common"
	"github.com/golang-jwt/jwt/v5"
)

func (cfg *apiConfig) GetNewJWT(user *common.User) {
	t := jwt.NewNumericDate(time.Now().UTC())
	if user.Expires_in_seconds == 0 {
		log.Println("No expiration given, defauilting to 60 seconds")
		user.Expires_in_seconds = 3600
	}

	expiretime := jwt.NewNumericDate(t.Add(time.Duration(user.Expires_in_seconds) * time.Second))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  t,
		ExpiresAt: expiretime,
		Subject:   strconv.Itoa(user.Id),
	})

	signedtoken, err := token.SignedString([]byte(cfg.secret_Key))
	if err != nil {
		log.Printf("Unable to generate a signed token %s\n", err)
		return
	}

	//test logged token
	log.Printf("JWT token generated %s\n", signedtoken)
	user.Token = signedtoken
}
