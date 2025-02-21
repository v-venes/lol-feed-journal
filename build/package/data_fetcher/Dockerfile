FROM golang:1.23

WORKDIR /app

COPY . .

RUN go build -o data_fetcher ./cmd/data_fetcher

RUN apt-get update && apt-get install -y cron

COPY ./deployments/cronjob /etc/cron.d/data_fetcher-cronjob

RUN chmod 0644 /etc/cron.d/data_fetcher-cronjob

RUN crontab /etc/cron.d/data_fetcher-cronjob

CMD ["cron", "-f"]