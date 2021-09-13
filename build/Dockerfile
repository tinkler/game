FROM golang:1.16.8
LABEL Author=tinkler
WORKDIR $GOPATH/src/sfs/cmd/server
ADD . $GOPATH/src/sfs
ENV GOPROXY=https://goproxy.cn
RUN go build .
EXPOSE 9301
ENTRYPOINT ["./server"]