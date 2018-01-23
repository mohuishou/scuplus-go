FROM golang:1.9.2-stretch
WORKDIR /go/src/github.com/mohuishou/scuplus-go
COPY . .
RUN go build -o /scuplus .
CMD ["/scuplus"]
