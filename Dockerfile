FROM golang:1.16
COPY . /go/src/github.com/polyverse/ropoly-cmd
WORKDIR /go/src/github.com/polyverse/ropoly-cmd
RUN go build
CMD bash