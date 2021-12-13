FROM golang:1.14-alpine as build

RUN apk add --no-cache git

WORKDIR /src

COPY . /src

RUN go mod vendor

RUN go build consumer.go

FROM alpine as runtime

COPY --from=build /src/consumer /app/consumer

CMD ["/app/consumer"]