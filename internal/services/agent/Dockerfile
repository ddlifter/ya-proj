FROM golang:latest

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o agent app/main.go

CMD ["./agent"]

#docker run --network=host 7aa7
