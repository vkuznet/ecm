package main

import (
	"encoding/json"
	"log"

	kvdb "github.com/vkuznet/ecm/kvdb"
)

// User represents our user attributes
type User struct {
	Name     string
	Password string
	Email    string
	Secret   string
}

// helper function to get user data
func getUser(name string) (User, error) {
	var user User
	data, err := DBStore.Get(name)
	if err != nil {
		log.Println("unable to get user", name, "error", err)
		return user, err
	}
	err = json.Unmarshal(data, &user)
	if err != nil {
		log.Printf("unable to unmarshal user '%s' data '%s', error %v", name, string(data), err)
		return user, err
	}
	return user, nil
}

// helper function to update user secret
func updateUser(name, secret string) {
	user, err := getUser(name)
	if err != nil {
		log.Println("unable to get user data", err)
		return
	}
	// delete existing record
	err = DBStore.Delete(name)
	if err != nil {
		log.Println("unable to delete user record", err)
		return
	}
	user.Secret = secret
	addUser(user)
}

// helper function to add user info
func addUser(user User) {
	data, err := json.Marshal(user)
	if err != nil {
		log.Println("unable to marshal user", err)
	}
	rec := kvdb.KVRecord{
		Key:   user.Name,
		Value: data,
	}
	err = DBStore.Add(rec)
	if err != nil {
		log.Println("unable to add user record to DBStore, error", err)
	}
}

// helper function to check if user exists in our DB
func userExist(name, password string) bool {
	user, err := getUser(name)
	if err != nil {
		return false
	}
	if password == "do not check" {
		return true
	} else if checkPasswordHash(password, user.Password) {
		return true
	}
	return false
}

// helper function to find user secret
func findUserSecret(name string) string {
	user, err := getUser(name)
	if err != nil {
		log.Println("unable to find user", err)
		return ""
	}
	return user.Secret
}
