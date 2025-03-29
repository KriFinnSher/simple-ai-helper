FROM golang:1.23

WORKDIR /support
COPY . .

RUN go mod tidy
RUN go build -o /build ./server/main.go

EXPOSE 8080
CMD ["/build"]