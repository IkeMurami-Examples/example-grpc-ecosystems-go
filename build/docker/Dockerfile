FROM golang AS builder

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GO111MODULE=on

RUN mkdir -p /go/src/github.com/z3f1r/grpc-ecosystems-example
ADD . /go/src/github.com/z3f1r/grpc-ecosystems-example

WORKDIR /go/src/github.com/z3f1r/grpc-ecosystems-example/

# RUN go test -mod vendor -v ./...
RUN go build -mod vendor -a -o /bin/start-example /go/src/github.com/z3f1r/grpc-ecosystems-example/cmd/exmaple/main.go


FROM alpine

RUN apk --no-cache add ca-certificates

WORKDIR /bin
COPY --from=builder /bin/start-example /bin/start-example

EXPOSE 8080
EXPOSE 7080

ENTRYPOINT ["/bin/start-example"]