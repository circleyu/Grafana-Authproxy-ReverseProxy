package main

import (
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
	if grafanaURL = os.Getenv("GRAFANA_URL"); grafanaURL == "" {
		grafanaURL = "http://localhost:3000/"
	}
}

func main() {
	remote, err := url.Parse(grafanaURL)
	if err != nil {
		panic(err)
	}

	proxy = httputil.NewSingleHostReverseProxy(remote)

	http.HandleFunc("/", mainHandler)

	err = http.ListenAndServe(":"+listenPort, nil)
	if err != nil {
		panic(err)
	}
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	if createAPI(w, r) {
		return
	}

	cookie := getCookie(r)

	logout(w, r, cookie)

	if isCookieExist(cookie) {
		proxyHandler(w, r)
		return
	}

	basicAuth(w, r)
}

func createAPI(w http.ResponseWriter, r *http.Request) bool {
	if r.Method == "POST" && strings.Contains(r.RequestURI, "/create/") {
		createHandler(w, r)
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
	m := r.URL.Query()

	userName, ok := m["user"]
	if !ok {
		// 認證失敗，提示 401 Unauthorized
		// Restricted 可以改成其他的值
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		// 401
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	r.SetBasicAuth(userName[0], userName[0])
	r.Header.Add("X-WEBAUTH-USER", userName[0])
	w.Header().Set("X-Frame-Options", "allow-from "+grafanaURL)
	setCookie(w, userName[0])
	proxy.ServeHTTP(w, r)
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Frame-Options", "allow-from "+grafanaURL)
	proxy.ServeHTTP(w, r)
}
