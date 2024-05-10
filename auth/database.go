package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/thanhpk/randstr"
)

type LoginDatabaseManager struct {
	db *sql.DB
}

type user struct {
	id    int64
	login string
	pass  string
	salt  string
}

type WrongCredentialsError struct {
	err error
}

func (e WrongCredentialsError) Error() string {
	return e.err.Error()
}

type LoginExistsError struct {
	err error
}

func (e LoginExistsError) Error() string {
	return e.err.Error()
}

func (dbm *LoginDatabaseManager) loginAccount(table, login, password string) (int64, error) {
	var userInstance user
	selectErr := dbm.db.QueryRow(fmt.Sprintf("SELECT * FROM %s WHERE login=$1", table), login).Scan(
		&userInstance.id,
		&userInstance.login,
		&userInstance.pass,
		&userInstance.salt,
	)
	if !errors.Is(selectErr, sql.ErrNoRows) {
		return 0, WrongCredentialsError{errors.New("wrong credentials")}
	} else if selectErr != nil {
		return 0, selectErr
	}
	hash := md5.Sum([]byte(password + userInstance.salt))
	if userInstance.pass != hex.EncodeToString(hash[:]) {
		return 0, WrongCredentialsError{errors.New("wrong credentials")}
	}
	return userInstance.id, nil
}
func (dbm *LoginDatabaseManager) createAccount(table, login, password string) (int64, error) {
	selectErr := dbm.db.QueryRow(fmt.Sprintf("SELECT login FROM %s WHERE login=$1", table), login).Scan()
	if selectErr == nil {
		return 0, LoginExistsError{errors.New("username already exists")}
	}
	if !errors.Is(selectErr, sql.ErrNoRows) {
		return 0, selectErr
	}
	salt := randstr.String(16)
	hash := md5.Sum([]byte(password + salt))
	hashedPass := hex.EncodeToString(hash[:])
	var id int64
	err := dbm.db.QueryRow(
		fmt.Sprintf("INSERT INTO %s (login, password, salt) VALUES ($1, $2, $3) RETURNING id", table),
		login, hashedPass, salt,
	).Scan(&id)
	return id, err
}

func (dbm *LoginDatabaseManager) LoginShop(login, password string) (int64, error) {
	return dbm.loginAccount("shops", login, password)
}
func (dbm *LoginDatabaseManager) CreateShop(login, password string) (int64, error) {
	return dbm.createAccount("shops", login, password)
}
func (dbm *LoginDatabaseManager) LoginUser(login, password string) (int64, error) {
	return dbm.loginAccount("users", login, password)
}
func (dbm *LoginDatabaseManager) CreateUser(login, password string) (int64, error) {
	return dbm.createAccount("users", login, password)
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
