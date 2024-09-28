FROM golang:1.22-alpine

WORKDIR /app
COPY . .

WORKDIR /app/db_migrations
RUN go build -o /weather-db-migration

CMD [ "/weather-db-migration" ]
