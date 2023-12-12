FROM golang:alpine AS application
EXPOSE 10015

FROM golang:alpine AS builder
RUN apk add --no-cache --update gcc g++

WORKDIR /fake_server
COPY . .

WORKDIR /fake_server/account_service
RUN go mod tidy
RUN CGO_ENABLED=1 go build

FROM application
WORKDIR /app

RUN CGO_ENABLED=1
COPY --from=builder /fake_server/account_service/fake_server /app/fake_server
COPY --from=builder /fake_server/account_service/core-config.yaml /app/core-config.yaml
COPY --from=builder /fake_server/account_service/server.crt /app/server.crt
COPY --from=builder /fake_server/account_service/server.key /app/server.key

RUN mkdir data
COPY --from=builder /fake_server/account_service/data/ /app/data/


CMD [ "/app/fake_server" ]