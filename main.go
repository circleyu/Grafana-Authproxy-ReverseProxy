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
	adminPassWord string
	grafanaURL    string
)

func init() {
	if adminPassWord = os.Getenv("ADMIN_PASSWORD"); adminPassWord == "" {
		adminPassWord = "admin"
	}
	if grafanaURL = os.Getenv("GRAFANA_URL"); grafanaURL == "" {
		grafanaURL = "http://localhost:3000/"
	}

	mapDashboardUID = make(map[string]string)
	mapDashboardUID["api-gateway"] = "AWSAPIGat"
	mapDashboardUID["autoscaling"] = "AWSAutoS"
	mapDashboardUID["billing"] = "AWSBillin"
	mapDashboardUID["cloudfront"] = "AWSCFront"
	mapDashboardUID["cloudwatch-browser"] = "AWSCWBrow"
	mapDashboardUID["ebs"] = "AWSEbs000"

	mapDashboardUID["ec2"] = "AWSEc2000"
	mapDashboardUID["ecs"] = "AWSECS000"
	mapDashboardUID["efs"] = "AWSEFS000"
	mapDashboardUID["elasticache-redis"] = "AWSECRedis"
	mapDashboardUID["elb-application-lb"] = "AWSAlb000"
	mapDashboardUID["elb-classic-lb"] = "AWSClb000"

	mapDashboardUID["emr-hadoop-2"] = "AWSEMRhadoop"
	mapDashboardUID["events"] = "AWSEvents"
	mapDashboardUID["lambda"] = "AWSLambda000"
	mapDashboardUID["logs"] = "AWSLogs00"
	mapDashboardUID["rds"] = "AWSRds000"
	mapDashboardUID["redshift"] = "kkvc4M0ik"

	mapDashboardUID["s3"] = "AWSS31iWk"
	mapDashboardUID["ses"] = "AWSSes000"
	mapDashboardUID["sns"] = "AWSSNS000"
	mapDashboardUID["sqs"] = "AWSSQS000"
	mapDashboardUID["vpn"] = "AWSVpn000"
}

func main() {
	remote, err := url.Parse(grafanaURL)
	if err != nil {
		panic(err)
	}

	proxy = httputil.NewSingleHostReverseProxy(remote)

	http.HandleFunc("/", mainHandler)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	if createAPI(w, r) {
		return
	}

	if addAPI(w, r) {
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

func addAPI(w http.ResponseWriter, r *http.Request) bool {
	if r.Method == "POST" && strings.Contains(r.RequestURI, "/add/") {
		addHandler(w, r)
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
	if !ok || userName[0] == "admin" {
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
	cookie := getCookie(r)
	r.SetBasicAuth(cookie.Value, cookie.Value)
	w.Header().Set("X-Frame-Options", "allow-from "+grafanaURL)
	proxy.ServeHTTP(w, r)
}
