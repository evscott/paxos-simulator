FROM golang:latest
LABEL maintainer="Eliot Scott <eliotvscott@gmail.com>"

WORKDIR /

COPY go.mod .
COPY artifacts /artifacts
COPY cmd /cmd
COPY /Multi //Multi

CMD ["go", "run", "cmd/main.go"]