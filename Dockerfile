# syntax=docker/dockerfile:1

FROM golang:alpine

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o /pd-payment ./src

EXPOSE 3006

CMD [ "/pd-payment" ]