FROM golang:alpine

COPY . /go/src/sample

RUN go install sample/server

EXPOSE 8080

CMD ["server"]
