package main

import (
	"log"
)

func (cfg *apiConfig) testLog(msg, value string) {
	log.Printf("%s : %s\n", msg, value)
}

func (cfg *apiConfig) isHashed(s string) bool {
	// Bcrypt hashes start with $2a$, $2b$, or $2y$
	return len(s) == 60 && (s[:4] == "$2a$" || s[:4] == "$2b$" || s[:4] == "$2y$")
}
