package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// SetCookie will set cookies in the user browser
func setCookie(w http.ResponseWriter, r *http.Request) {

	// Set expiration time
	expiration := time.Now()
	expiration = time.Now().Add(time.Second * time.Duration(15))

	// Generate cookie
	cookie := http.Cookie{Name: "grafana_session", Value: "rain", Path: "/", Expires: expiration}
	// Set cookie
	http.SetCookie(w, &cookie)
}

func clearCookie(w http.ResponseWriter, cookie *http.Cookie) {
	if cookie == nil {
		return
	}
	cookie.Value = ""
	// Set cookie
	http.SetCookie(w, cookie)
	log.Println(fmt.Sprintf("clear cookie : %s", cookie))
}

func getCookie(w http.ResponseWriter, r *http.Request) *http.Cookie {
	if cookie, err := r.Cookie("grafana_session"); err == nil {
		return cookie
	}
	return nil
}

func isCookieExist(cookie *http.Cookie) bool {
	if cookie == nil {
		return false
	}
	return cookie.Value != ""
}
