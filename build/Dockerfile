FROM golang:1.20 as builder

ARG GITHUB_USERNAME
ARG GITHUB_TOKEN

ENV HOME /app
ENV CGO_ENABLED 0
ENV GOOS linux

WORKDIR /app
COPY go.mod go.sum ./
RUN echo "machine github.com login $GITHUB_USERNAME password $GITHUB_TOKEN" > $HOME/.netrc
RUN go env -w GOPRIVATE=github.com/vediagames/*
RUN go mod download
COPY . .

RUN go build -a -installsuffix cgo -o server .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /

COPY --from=builder /app/server .

CMD [ "./server" ]
