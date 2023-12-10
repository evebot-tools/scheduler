FROM golang:1.21 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-s -w" -o #APP# ./cmd

FROM gcr.io/distroless/static:nonroot as app

LABEL org.opencontainers.image.source="https://github.com/OWNER/REPO"

USER nonroot:nonroot

COPY --from=builder --chown=nonroot:nonroot /app/#APP# /#APP#

ENTRYPOINT ["/#APP#"]