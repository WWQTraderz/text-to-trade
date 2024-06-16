FROM --platform=$TARGETARCH golang:1.22-alpine

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY go.mod go.sum ./
COPY ./ .

RUN go mod download
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /app/api /app/cmd/api

ENTRYPOINT [ "/app/api" ]