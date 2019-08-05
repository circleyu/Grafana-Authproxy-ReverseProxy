FROM golang:latest AS buildStage
WORKDIR /go/src/ReverseProxy
COPY . .
RUN CGO_ENABLED=0 go build

FROM scratch
WORKDIR /app
COPY --from=buildStage /go/src/ReverseProxy/ReverseProxy .
COPY --from=buildStage /go/src/ReverseProxy/dashboard ./dashboard/
EXPOSE 8080
ENTRYPOINT ["/app/ReverseProxy"]
