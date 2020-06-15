FROM golang:1.14 as builder

RUN mkdir -p /build
WORKDIR /build

ADD go.sum .
ADD go.mod .

RUN go mod download

ADD . .

RUN go build -o sentrygatewayd ./cmd

FROM ubuntu:19.10 as runner

COPY --from=builder /build/sentrygatewayd /bin/sentrygatewayd
ADD config.yaml .

EXPOSE 80 80

CMD ["/bin/sentrygatewayd"]
