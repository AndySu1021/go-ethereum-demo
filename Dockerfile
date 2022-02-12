FROM golang:1.16 as build

WORKDIR /ethereum

COPY . .

RUN go build -o api main.go
RUN go build -o scan scan.go
RUN go build -o confirm confirm.go

FROM gcr.io/distroless/base

WORKDIR /ethereum

COPY --from=build /ethereum /ethereum

CMD ["./api"]