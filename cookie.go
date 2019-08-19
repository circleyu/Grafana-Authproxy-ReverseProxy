package main

import (
	"net/http"
	"time"
)

// SetCookie will set cookies in the user browser
func setCookie(w http.ResponseWriter, v string) {

	// Set expiration time
	expiration := time.Now()
	expiration = time.Now().Add(time.Hour * time.Duration(1))

	// Generate cookie
	cookie := http.Cookie{Name: "grafana_session", Value: v, Path: "/", Expires: expiration}
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
}

func getCookie(r *http.Request) *http.Cookie {
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
