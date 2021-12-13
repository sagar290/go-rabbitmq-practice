FROM golang:1.14-alpine as build

RUN apk add --no-cache git

WORKDIR /src

COPY . /src

RUN go mod vendor

RUN go build producer.go

FROM alpine as runtime

COPY --from=build /src/producer /app/producer

CMD ["/app/producer"]