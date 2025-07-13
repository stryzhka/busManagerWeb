FROM golang:latest
WORKDIR /go/src/app

COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY . .

RUN swag init -g cmd/server/main.go

EXPOSE 8080

CMD ["go", "run", "cmd/server/main.go"]

