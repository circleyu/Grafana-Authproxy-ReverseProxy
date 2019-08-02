package main

import (
	"bytes"
	"encoding/base64"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

type viewFunc func(http.ResponseWriter, *http.Request)

var proxy *httputil.ReverseProxy

func main() {
	remote, err := url.Parse("http://localhost:3000/")
	if err != nil {
		panic(err)
	}

	proxy = httputil.NewSingleHostReverseProxy(remote)

	// 塞入要pass的handler
	http.HandleFunc("/", mainHandler)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	if api(w, r) {
		return
	}

	logout(w, r)
	cookie := getCookie(w, r)

	if isCookieExist(cookie) {
		proxyHandler(w, r)
		return
	}

	basicAuth(w, r)
}

func api(w http.ResponseWriter, r *http.Request) bool {
	if strings.Contains(r.RequestURI, "/api/") {
		r.Header.Del("X-WEBAUTH-USER")
		proxyHandler(w, r)
		return true
	}
	return false
}

func logout(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.RequestURI, "/logout") {
		log.Println("logout")
		delCookie(w, r)
	}
}

func basicAuth(w http.ResponseWriter, r *http.Request) {

	user := []byte("admin")
	passwd := []byte("admin")

	basicAuthPrefix := "Basic "

	// 獲取 request header
	auth := r.Header.Get("Authorization")
	// 如果是 http basic auth
	if strings.HasPrefix(auth, basicAuthPrefix) {
		// 解析認證訊息
		payload, err := base64.StdEncoding.DecodeString(
			auth[len(basicAuthPrefix):],
		)
		if err == nil {
			pair := bytes.SplitN(payload, []byte(":"), 2)
			if len(pair) == 2 && bytes.Equal(pair[0], user) &&
				bytes.Equal(pair[1], passwd) {
				// 執行函式
				setCookie(w, r)
				proxyHandler(w, r)
				return
			}
		}
	}

	// 認證失敗，提示 401 Unauthorized
	// Restricted 可以改成其他的值
	w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
	// 401
	w.WriteHeader(http.StatusUnauthorized)
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Frame-Options", "allow-from http://localhost:3000/")
	proxy.ServeHTTP(w, r)
}

// SetCookie will set cookies in the user browser
func setCookie(w http.ResponseWriter, r *http.Request) {

	// Set expiration time
	expiration := time.Now()
	expiration = expiration.AddDate(1, 0, 0)

	// Generate cookie
	cookie := http.Cookie{Name: "username", Value: "rain", Expires: expiration}

	// Set cookie
	http.SetCookie(w, &cookie)
}

func delCookie(w http.ResponseWriter, r *http.Request) {
	// Generate cookie
	cookie := getCookie(w, r)
	cookie.MaxAge = -1
	// Set cookie
	http.SetCookie(w, cookie)
}

func getCookie(w http.ResponseWriter, r *http.Request) *http.Cookie {
	for _, cookie := range r.Cookies() {
		log.Println(cookie)
	}
	if cookie, err := r.Cookie("username"); err == nil {

		//log.Println(cookie)
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
