# grafana-auth-proxy
透過Reverse Proxy以及URL中的Query內容，來實現快速登入Grafana的功能

## 執行
這邊是透過docker來執行的，在[docker-compose.yml](https://github.com/grandtechcloud/Grafana-Authproxy-ReverseProxy/blob/master/docker-compose.yml)裡面有環境變數需要設定

```json
environment:
  - LISTEN_PORT=8080
  - ADMIN_PASSWORD=admin
  - GRAFANA_URL=http://grafana:3000/
```

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

## dashboard
目前預設只會建置EC2的dashboard<br>
其餘的請手動在Grafana建置，直接在Grafana導入[dashboard](https://github.com/grandtechcloud/Grafana-Authproxy-ReverseProxy/tree/master/dashboard)目錄下的json檔案

## 實作參考
[這邊](https://www.annhe.net/article-3551.html)
