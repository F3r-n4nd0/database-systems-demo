FROM golang:latest as build

WORKDIR $GOPATH/src/web-service-kudos

ADD . .

# USE WITH VGO
#RUN go version && go get -u -v golang.org/x/vgo
#ENV CGO_ENABLED=0
#RUN vgo install ./...

RUN go get -d -v ./...

EXPOSE 8080

CMD ["go", "run", "main.go"]
