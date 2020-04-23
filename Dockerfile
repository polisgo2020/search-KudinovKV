FROM golang:1.13 AS builder
WORKDIR /usr/src

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -installsuffix cgo -o search-KudinovKV main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /usr/app
COPY --from=builder /usr/src/search-KudinovKV .
COPY --from=builder /usr/src/web ./web
ENTRYPOINT ["/usr/app/search-KudinovKV"]
CMD ["search" , "database"]