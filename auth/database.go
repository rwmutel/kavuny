package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/thanhpk/randstr"
	"net/http"
	"slices"
)

type UserType string

var userTypes = []UserType{"user", "shop"}

type LoginDatabaseManager struct {
	db *sql.DB
}

type user struct {
	id    int64
	login string
	pass  string
	salt  string
	user  UserType
}

func (dbm *LoginDatabaseManager) CheckUserType(userType UserType) bool {
	return slices.Contains(userTypes, userType)
}
func (dbm *LoginDatabaseManager) loginAccount(login, password string) (int64, UserType, *HttpError) {
	var userInstance user
	selectErr := dbm.db.QueryRow("SELECT * FROM users WHERE login=$1", login).Scan(
		&userInstance.id,
		&userInstance.login,
		&userInstance.pass,
		&userInstance.salt,
		&userInstance.user,
	)
	if errors.Is(selectErr, sql.ErrNoRows) {
		return 0, "", NewHttpError(nil, "Wrong credentials", http.StatusForbidden)
	} else if selectErr != nil {
		return 0, "", NewHttpError(selectErr, "Failed to select user", http.StatusServiceUnavailable)
	}
	hash := md5.Sum([]byte(password + userInstance.salt))
	if userInstance.pass != hex.EncodeToString(hash[:]) {
		return 0, "", NewHttpError(nil, "Wrong credentials", http.StatusForbidden)
	}
	return userInstance.id, userInstance.user, nil
}
func (dbm *LoginDatabaseManager) CreateAccount(login, password string, user UserType) (int64, *HttpError) {
	var userCount int
	selectErr := dbm.db.QueryRow("SELECT count(*) FROM users WHERE login=$1", login).Scan(&userCount)
	if selectErr != nil {
		return 0, NewHttpError(selectErr, "Failed to select user", http.StatusServiceUnavailable)
	} else if userCount > 0 {
		return 0, NewHttpError(nil, "Username already exists", http.StatusForbidden)
	}
	salt := randstr.String(16)
	hash := md5.Sum([]byte(password + salt))
	hashedPass := hex.EncodeToString(hash[:])
	var id int64
	err := dbm.db.QueryRow(
		"INSERT INTO users (login, password, salt, user_type) VALUES ($1, $2, $3, $4) RETURNING id",
		login, hashedPass, salt, user,
	).Scan(&id)
	if err != nil {
		return 0, NewHttpError(err, "Failed to insert user", http.StatusServiceUnavailable)
	}
	return id, nil
}

func (dbm *LoginDatabaseManager) Initialize(user, password, host, port string) (err error) {
	connStr := fmt.Sprintf("sslmode=disable user=%s password=%s host=%s port=%s",
		user, password, host, port)
	dbm.db, err = sql.Open("postgres", connStr)
	return
}

func (dbm *LoginDatabaseManager) Close() error {
	return dbm.db.Close()
}
