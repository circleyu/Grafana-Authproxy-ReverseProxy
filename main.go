package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type viewFunc func(http.ResponseWriter, *http.Request)

var proxy *httputil.ReverseProxy

func main() {
	remote, err := url.Parse("http://grafana:3000/")
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

	cookie := getCookie(w, r)

	logout(w, r, cookie)

	log.Println(fmt.Sprintf("get cookie : %s", cookie))

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

func logout(w http.ResponseWriter, r *http.Request, cookie *http.Cookie) {
	if strings.Contains(r.RequestURI, "logout") {
		clearCookie(w, cookie)
		r.Header.Del("Authorization")
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
				r.Header.Add("X-WEBAUTH-USER", string(user))
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
	w.Header().Set("X-Frame-Options", "allow-from http://grafana:3000/")
	proxy.ServeHTTP(w, r)
}