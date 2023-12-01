FROM golang:alpine AS application
EXPOSE 10015

FROM golang:alpine AS builder
RUN apk add --no-cache --update gcc g++

WORKDIR /fake_server
COPY . .

RUN go mod tidy
RUN CGO_ENABLED=1 go build

FROM application
WORKDIR /app

RUN CGO_ENABLED=1
COPY --from=builder /fake_server/fake_server /app/fake_server
RUN mkdir data

CMD [ "/app/fake_server" ]