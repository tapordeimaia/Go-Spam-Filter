FROM golang:1.25-alpine

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY main.go ./
COPY model.json ./

RUN go build -o spam-api

EXPOSE 8080

CMD ["./spam-api"]