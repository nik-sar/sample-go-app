FROM golang:1.22-alpine

WORKDIR /app
COPY main.go go.sum go.mod ./
RUN CGO_ENABLED=0 GOOS=linux go build -o ./sample-go-app
EXPOSE 8080
CMD ["./sample-go-app"]