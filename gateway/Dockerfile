FROM golang:latest

WORKDIR /app

COPY gateway/go.mod .
COPY gateway/go.sum .

RUN go mod download

COPY gateway .

RUN go build -o gateway cmd/app/main.go

CMD ["./gateway"]

#docker run --network=host server

