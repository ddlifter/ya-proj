FROM golang:latest

WORKDIR /app

COPY agent/go.mod .
COPY agent/go.sum .

RUN go mod download

COPY agent .

RUN go build -o agent app/main.go

CMD ["./agent"]

#docker run --network=host 7aa7
