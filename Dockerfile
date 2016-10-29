FROM golang:alpine

ADD . /go/src/github.com/mikelong/usersvc

RUN go install github.com/mikelong/usersvc/cmd/usersvc

ENTRYPOINT /go/bin/usersvc

EXPOSE 8080
