package handlers

import (
	"encoding/json"
	"github.com/gofrs/uuid"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	u = "stratton"
	p = ".dev"
	path = "/tmp/"
	cookieName = "strattonDevSession"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	var credentials Credentials

	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if credentials.Username != u || credentials.Password != p {
		log.Println("Credentials were incorrect")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token, err := uuid.NewV4()

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = os.Create(path + token.String())

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name: cookieName,
		Value: token.String(),
		Expires: time.Now().Add(24 * time.Hour),
	})
}


func Authenticated(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(cookieName)
	if err != nil {
		if err == http.ErrNoCookie {
			log.Println("Could not find session cookie on request")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, err := os.Stat(path + c.Value); os.IsNotExist(err) {
		log.Println("Session did not exist")
		w.WriteHeader(http.StatusUnauthorized)
	} else {
		_, _ = w.Write([]byte("Welcome\n"))
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(cookieName)
	if err != nil {
		if err == http.ErrNoCookie {
			log.Println("Returning success on logout request without session cookie")
			return
		}

		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = os.Remove(path + c.Value)

	if err != nil {
		if os.IsNotExist(err) {
			return
		}

		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
