package main

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
)

const saltSize = 64

type authStruct struct {
	Username string
	Password string
}

func parseReq(request *http.Request) authStruct {
	decoder := json.NewDecoder(request.Body)
	var n authStruct
	err := decoder.Decode(&n)
	if err != nil {
		panic(err)
	}
	return n
}

func generateSalt(saltSize int) []byte {
	var salt = make([]byte, saltSize)
	_, err := rand.Read(salt[:])
	if err != nil {
		panic(err)
	}
	return salt
}

func saltAndHash(password string, salt []byte) string {
	saltedPassword := append([]byte(password), salt...)

	sha512Hasher := sha512.New()
	sha512Hasher.Write(saltedPassword)

	saltedHash := sha512Hasher.Sum(nil)
	return hex.EncodeToString(saltedHash)
}

func main() {
	fmt.Println("Start!")

	mux := http.NewServeMux()

	mux.HandleFunc("POST /signup", func(w http.ResponseWriter, r *http.Request) {
		decodedBody := parseReq(r)
		//TO DO: check DB for if username is unique
		userSalt := generateSalt(saltSize)
		passwordHash := saltAndHash(decodedBody.Password, userSalt)
		//TO DO: Add user to DB (username,passHash,salt,maybe id?)
		fmt.Fprintf(w, "Hello! The hash is %s.", passwordHash)
	})

	mux.HandleFunc("GET /login", func(w http.ResponseWriter, r *http.Request) {
		decodedBody := parseReq(r)
		//TO DO: check DB for username
		//TO DO: retrieve salt from DB
		//TO DO: retrieve password hash from DB
		userSalt := generateSalt(saltSize) //Remove once DB works
		passwordHash := saltAndHash(decodedBody.Password, userSalt)
		//TO DO: compare hashes
		fmt.Fprintf(w, "Hello! The hash is %s.", passwordHash)
		fmt.Fprintf(w, "Hello, %s! We're such great friends!", decodedBody.Username)
	})

	mux.HandleFunc("DELETE /delete", func(w http.ResponseWriter, r *http.Request) {
		decodedBody := parseReq(r)
		//TO DO: check DB for username
		//TO DO: retrieve salt from DB
		//TO DO: retrieve password hash from DB
		userSalt := generateSalt(saltSize) //Remove once DB works
		passwordHash := saltAndHash(decodedBody.Password, userSalt)
		//TO DO: Add user to DB (username,passHash,salt,maybe id?)
		fmt.Fprintf(w, "The hash is %s.", passwordHash)
		fmt.Fprintf(w, "Goodbye %s. I'm so sorry to see you go. :(", decodedBody.Username)
	})

	if err := http.ListenAndServe("localhost:8080", mux); err != nil {
		fmt.Println(err.Error())
	}
}
