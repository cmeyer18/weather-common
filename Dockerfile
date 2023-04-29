# specify the base image to  be used for the application, alpine or ubuntu
FROM golang:1.19-alpine

# create a working directory inside the image
WORKDIR /app

COPY . .
RUN go install ./...

# download Go modules and dependencies
RUN go mod download

# download Go vendor
RUN go mod vendor

# compile application
RUN go build ./...
