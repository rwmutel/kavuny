package main

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/gocarina/gocsv"
	"github.com/thanhpk/randstr"
	"os"
)

type BareUser struct {
	Login    string `csv:"login"`
	Password string `csv:"password"`
	UserType string `csv:"user_type"`
}

type User struct {
	Id       int    `csv:"id"`
	Login    string `csv:"login"`
	Password string `csv:"password"`
	Salt     string `csv:"salt"`
	UserType string `csv:"user_type"`
}

func main() {
	inputFile, _ := os.Open("user_creds.csv")
	defer inputFile.Close()

	var users []BareUser
	gocsv.UnmarshalFile(inputFile, &users)

	var result []User
	for idx, userObj := range users {
		salt := randstr.String(16)
		hash := md5.Sum([]byte(userObj.Password + salt))
		hashedPass := hex.EncodeToString(hash[:])
		result = append(result, User{idx + 1, userObj.Login,
			hashedPass, salt, userObj.UserType})
	}

	out, _ := os.Create("users.csv")
	gocsv.MarshalFile(result, out)
}
