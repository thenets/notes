FROM docker.io/golang:latest as builder

ENV GOPATH=/go

WORKDIR /go/pkg/mod/github.com/thenets/notes/

COPY . .

RUN set -ex \
    && go mod download \
    && CGO_ENABLED=0 GOOS=linux go build -o /notes

# ---
# Final image
# ---
FROM docker.io/redhat/ubi9:9.4-1214.1726694543

WORKDIR /app

COPY --from=builder /notes /app/notes
COPY --from=builder /go/pkg/mod/github.com/thenets/notes/static /app/static

USER nonroot:nonroot

EXPOSE 8080

ENV PORT=8080

ENTRYPOINT ["/app/notes"]
