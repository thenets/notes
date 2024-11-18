FROM registry.access.redhat.com/ubi9/go-toolset:1.21.13-2.1729776560 as builder

ENV GOPATH=/go

USER root

WORKDIR /go/pkg/mod/github.com/thenets/notes/

COPY . .

RUN set -ex \
    && go mod download \
    && CGO_ENABLED=0 GOOS=linux go build -o /notes

# ---
# Final image
# ---
FROM docker.io/redhat/ubi9:9.5

WORKDIR /app

COPY --from=builder /notes /app/notes
COPY --from=builder /go/pkg/mod/github.com/thenets/notes/static /app/static

RUN set -x \
    && useradd -m -s /sbin/nologin -u 1000 notes \
    && chown -R notes:notes /app \
    && chmod -R 755 /app

USER notes

EXPOSE 8080

ENV PORT=8080

ENTRYPOINT ["/app/notes"]
