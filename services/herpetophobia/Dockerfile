FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod /app/
COPY go.sum /app/
RUN go mod download
COPY Makefile /app/
COPY service/backend/herpetophobia /app/
RUN go build -o srv
CMD ["./srv"]