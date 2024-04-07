FROM golang:1.22-alpine

WORKDIR /app
COPY . .

RUN go build -o /weather-db-migration

CMD [ "/weather-db-migration" ]
