FROM golang:latest AS builder
WORKDIR /go/src/app
COPY main.go go.mod go.sum ./
RUN go mod tidy
RUN go build -o /usr/bin/freenom-webhook .

###############################################

FROM alpine:latest as prod

RUN apk --no-cache add ca-certificates
RUN apk --no-cache add libc6-compat

COPY --from=builder /usr/bin/freenom-webhook /usr/bin/freenom-webhook

ENTRYPOINT ["/usr/bin/freenom-webhook"]