package main

import (
	"bytes"
	"encoding/base64"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type viewFunc func(http.ResponseWriter, *http.Request)

func main() {
	remote, err := url.Parse("http://localhost:3000/")
	if err != nil {
		panic(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)

	// 塞入要pass的handler
	http.HandleFunc("/", basicAuth(proxyHandler(proxy)))

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func basicAuth(f viewFunc) viewFunc {

	user := []byte("admin")
	passwd := []byte("admin")

	return func(w http.ResponseWriter, r *http.Request) {
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
					f(w, r)
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
}

func proxyHandler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL)
		w.Header().Set("X-Frame-Options", "allow-from http://localhost:3000/")
		p.ServeHTTP(w, r)
	}
}
