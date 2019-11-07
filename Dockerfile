FROM golang:1.13.4-stretch
WORKDIR /go/src/github.com/mohuishou/scuplus-go
COPY . .
RUN go build -o /scuplus .
CMD ["/scuplus"]
