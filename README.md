## 介紹
透過Reverse Proxy以及URL中的Query內容，來實現快速登入Grafana的功能<br>
```
http://localhost:8080/dashboards?user=使用者名稱
```

## 執行
這邊是透過docker來執行的，在[docker-compose.yml](https://github.com/grandtechcloud/Grafana-Authproxy-ReverseProxy/blob/master/docker-compose.yml)裡面有環境變數需要設定

```json
environment:
  - ADMIN_PASSWORD=admin
  - GRAFANA_URL=http://grafana:3000/
```
## 在新增User之前
需要在AWS帳戶新增一個只有AWS CloudWatch讀取權限的使用者

## 新增User
提供了API接口，只需要將下方的資料結構POST到create這個API即可<br>
ex:
```
http://localhost:8080/create/
```
 
```json
{
	"userName": "",
	"email": "",
	"accessKey": "",
	"secretKey": ""
}
```
目前會預設建置EC2的dashboard

## 新增dashboard
提供了API接口，只需要將下方的資料結構POST到add這個API即可<br>
ex:
```
http://localhost:8080/add/
```

```json
{
	"userName": "",
	"dashboard": [""]
}
```
dashboard請填入[dashboard](https://github.com/grandtechcloud/Grafana-Authproxy-ReverseProxy/tree/master/dashboard)目錄下的json檔案名稱，如果填入all則會全部都建置。

## Grafana admin登入
反向代理伺服器有擋Admin的登入，以防被他人登入後串改<br>
Admin可以透過Grafana的URL來登入並管理內容
