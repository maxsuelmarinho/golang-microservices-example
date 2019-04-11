FROM golang:alpine as builder

RUN apk update && apk add git && apk add ca-certificates
RUN adduser -D -g '' appuser
COPY . /app/accountservice/
WORKDIR /app/accountservice/
RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/accountservice-linux-amd64
RUN chown -R appuser:appuser /go

FROM iron/base
RUN adduser -D -g '' appuser
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /go/bin/accountservice-linux-amd64 /go/bin/accountservice-linux-amd64

USER appuser
ENTRYPOINT ["/go/bin/accountservice-linux-amd64"]

