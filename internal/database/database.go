package database

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"

	"github.com/calamityesp/chirpy/common"
	"golang.org/x/crypto/bcrypt"
)

type DB struct {
	path string
	mu   *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]common.Chirp `json:"chirps"`
	Users  map[int]common.User  `json:"user"`
}

func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mu:   &sync.RWMutex{},
	}
	err := db.ensureDB()
	return db, err
}

func (db *DB) GetUsers() ([]common.User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	users := make([]common.User, 0, len(dbStructure.Users))
	for _, user := range dbStructure.Users {
		users = append(users, user)
	}

	return users, nil
}

func (db *DB) GetUserByEmail(email string) (common.User, error) {
	var fUser common.User
	var emptyUser common.User

	dbStructure, err := db.loadDB()
	if err != nil {
		return common.User{}, err
	}

	for _, user := range dbStructure.Users {
		if user.Email == email {
			fUser = user
			break
		}
	}

	if fUser != emptyUser {
		return fUser, nil
	}
	return emptyUser, nil
}

func (db *DB) GetUserByID(id int) (common.User, error) {
	var fUser common.User
	var emptyUser common.User

	dbStructure, err := db.loadDB()
	if err != nil {
		return common.User{}, err
	}

	for _, user := range dbStructure.Users {
		if user.Id == id {
			fUser = user
			break
		}
	}

	if fUser != emptyUser {
		return fUser, nil
	}
	return emptyUser, nil
}

func (db *DB) CreateUser(body common.User) (common.User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return common.User{}, err
	}

	// has the password value from body
	hashedPassword, err := db.convertPasswordToHash(body.Password)
	if err != nil {
		return common.User{}, err
	}

	id := len(dbStructure.Users) + 1
	user := common.User{
		Id:       id,
		Email:    body.Email,
		Password: hashedPassword,
	}
	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return common.User{}, err
	}

	return user, nil
}

func (db *DB) UpdateUser(update common.User) (common.User, error) {
	found := false

	DBStructure, err := db.loadDB()
	if err != nil {
	}

	for key, user := range DBStructure.Users {
		if key == update.Id {
			found = true
			user.Id = update.Id
			user.Email = update.Email
			user.Password = update.Password
			DBStructure.Users[key] = user
			break
		}
	}

	if found == false {
		return update, errors.New("User not found")
	}

	// delete then rewqrite the database
	db.deleteDatabase()
	db.writeDB(DBStructure)

	log.Printf("Updated User: id-%d, email-%s", update.Id, update.Email)
	return update, nil
}

func (db *DB) CreateChirp(body string) (common.Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return common.Chirp{}, err
	}

	id := len(dbStructure.Chirps) + 1
	chirp := common.Chirp{
		ID:   id,
		Body: body,
	}
	dbStructure.Chirps[id] = chirp

	err = db.writeDB(dbStructure)
	if err != nil {
		return common.Chirp{}, err
	}

	return chirp, nil
}

func (db *DB) GetChirps() ([]common.Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]common.Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}

func (db *DB) GetChirpById(id int) (common.Chirp, error) {
	var fChirp common.Chirp
	var emptyChirp common.Chirp

	dbStructure, err := db.loadDB()
	if err != nil {
		return common.Chirp{}, err
	}

	for _, chirp := range dbStructure.Chirps {
		if chirp.ID == id {
			fChirp = chirp
			break
		}
	}

	if fChirp != emptyChirp {
		return fChirp, nil
	}
	return emptyChirp, nil
}

func (db *DB) createDB() error {
	dbStructure := DBStructure{
		Chirps: map[int]common.Chirp{},
		Users:  map[int]common.User{},
	}
	return db.writeDB(dbStructure)
}

func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return db.createDB()
	}
	return err
}

func (db *DB) loadDB() (DBStructure, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	dbStructure := DBStructure{}
	dat, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return dbStructure, err
	}
	err = json.Unmarshal(dat, &dbStructure)
	if err != nil {
		return dbStructure, err
	}

	return dbStructure, nil
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	dat, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, dat, 0600)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) deleteDatabase() error {
	db.mu.Lock()
	defer db.mu.Unlock()

	err := os.Remove(db.path)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) convertPasswordToHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	hashToString := string(hash)
	return hashToString, nil
}
