package main

import (
	"bytes"
	"encoding/base64"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

type viewFunc func(http.ResponseWriter, *http.Request)

var (
	proxy         *httputil.ReverseProxy
	listenPort    string
	adminPassWord string
	grafanaURL    string
)

func init() {
	if listenPort = os.Getenv("LISTEN_PORT"); listenPort == "" {
		listenPort = "8080"
	}
	if adminPassWord = os.Getenv("ADMIN_PASSWORD"); adminPassWord == "" {
		adminPassWord = "admin"
	}
	if grafanaURL = os.Getenv("GRAFANA_URL"); adminPassWord == "" {
		grafanaURL = "http://localhost:3000/"
	}
}

func main() {
	remote, err := url.Parse(grafanaURL)
	if err != nil {
		panic(err)
	}

	proxy = httputil.NewSingleHostReverseProxy(remote)

	// 塞入要pass的handler
	http.HandleFunc("/", mainHandler)

	err = http.ListenAndServe(":"+listenPort, nil)
	if err != nil {
		panic(err)
	}
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	if grafana(w, r) {
		return
	}

	if api(w, r) {
		return
	}

	cookie := getCookie(w, r)

	logout(w, r, cookie)

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

func grafana(w http.ResponseWriter, r *http.Request) bool {
	if strings.Contains(r.RequestURI, "/grafana/") {
		r.Header.Del("X-WEBAUTH-USER")
		r.Header.Add("Authorization", "Basic YWRtaW46YWRtaW4=")
		r.Header.Add("X-WEBAUTH-USER", "admin")
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
	w.Header().Set("X-Frame-Options", "allow-from http://localhost:3000/")
	proxy.ServeHTTP(w, r)
}
