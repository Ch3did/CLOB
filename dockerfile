FROM golang:1.24

WORKDIR /app

COPY go.mod go.sum ./

COPY . .

RUN go build -o app ./api

CMD ["./app"]