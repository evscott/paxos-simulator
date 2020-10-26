FROM golang:latest

LABEL maintainer="Eliot Scott <eliotvscott@gmail.com>"

WORKDIR /

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o basic-paxos .

EXPOSE 9090

CMD ["./basic-paxos"]