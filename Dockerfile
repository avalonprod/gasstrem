 16 lines (9 sloc) 154 Bytes
FROM golang:1.19

EXPOSE 8000

WORKDIR /app

COPY go.mod ./
# COPY go.sum ./

RUN go mod download

COPY . .

RUN go build src/cmd/app/main.go 

CMD ["./main"]
