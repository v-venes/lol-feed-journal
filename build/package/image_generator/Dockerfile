FROM golang:1.23

WORKDIR /app

COPY . .

RUN go build -o image_generator ./cmd/image_generator

CMD ["./image_generator"]