FROM golang:latest AS buildStage
WORKDIR /ReverseProxy
COPY go.mod go.sum ./
RUN  go mod download
COPY . .
RUN CGO_ENABLED=0 go build

FROM scratch
WORKDIR /app
COPY --from=buildStage /ReverseProxy/ReverseProxy .
COPY --from=buildStage /ReverseProxy/dashboard ./dashboard/
EXPOSE 8080
ENTRYPOINT ["/app/ReverseProxy"]
