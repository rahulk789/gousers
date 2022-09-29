#Build stage 1
FROM golang:1.18 AS builder
WORKDIR /gousers
COPY go.mod /gousers/go.mod
COPY main.go /gousers/main.go
RUN go mod tidy
RUN go build -o main main.go

#Build stage 2
FROM alpine
WORKDIR /gousers
COPY --from=builder /gousers/main . 
EXPOSE 8080
CMD ["/gousers/main"]

