FROM golang:1.21-alpine

WORKDIR /app
COPY . .

RUN go build -o /weather-db-migration

CMD [ "/weather-db-migration" ]
