FROM golang:1.23

WORKDIR /app

COPY . .

RUN go build -o disc_bot ./cmd/disc_bot

CMD ["./disc_bot"]