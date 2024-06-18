package main

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/calamityesp/chirpy/common"
	"github.com/golang-jwt/jwt/v5"
)

func (cfg *apiConfig) GetNewJWT(user *common.User) {
	defaultExpiration := 60 * 60
	if user.Expires_in_seconds == 0 || user.Expires_in_seconds > defaultExpiration {
		log.Println("No expiration given, defauilting to 60 seconds")
		user.Expires_in_seconds = defaultExpiration
	}

	log.Printf("ExpirationTIme: %d", user.Expires_in_seconds)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Duration(user.Expires_in_seconds) * time.Second)),
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

	if user.RefreshToken == "" {
		cfg.GetRefreshToken(user)
	}
}

func (cfg *apiConfig) GetRefreshToken(user *common.User) {
	refreshTokenExpire := 60 * 60 * 24 * 60 // expire in 60 days
	numBytes := 32
	bSlice := make([]byte, numBytes)
	_, err := rand.Read(bSlice)
	if err != nil {
		log.Printf("Error generating random value: %s", err.Error())
		return
	}
	user.RefreshToken = hex.EncodeToString(bSlice)
	user.Refresh_token_expire_time = time.Now().Add(time.Duration(refreshTokenExpire) * time.Second)
	cfg.DB.UpdateUserRefreshToken(user)
	cfg.testLog("refreshToken: ", user.RefreshToken)
}

func (cfg *apiConfig) ValidateRefreshToken(user *common.User) bool {
	timeCompare := time.Now().UTC().Compare(user.Refresh_token_expire_time)
	cfg.testLog("Validate Refresh Token", strconv.Itoa(timeCompare))
	cfg.testLog("Validate Refresh Token", user.Refresh_token_expire_time.String())
	if timeCompare > 0 {
		log.Println("Current time is before given time")
		return false
	}
	return true
}

func (cfg *apiConfig) validateJwtToken(authHeader string) (common.User, bool) {
	if authHeader == "" {
		return common.User{}, false
	}

	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "Bearer" {
		return common.User{}, false
	}
	tokenString := splitAuth[1]
	log.Printf("tokenstring: %s\n", tokenString)

	secret := cfg.secret_Key
	claim := jwt.RegisteredClaims{}

	cfg.testLog("validating ", tokenString)
	token, err := jwt.ParseWithClaims(tokenString, &claim, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		log.Printf("You fucked up %s\n", err)
		return common.User{}, false
	}

	// validate token
	if !token.Valid {
		return common.User{}, false
	}

	subject, err := token.Claims.GetSubject()
	if err != nil {
		return common.User{}, false
	}

	// convert userId to int
	userId, err := strconv.Atoi(subject)
	if err != nil {
		return common.User{}, false
	}

	user, err := cfg.DB.GetUserByID(userId)
	return user, true
}
